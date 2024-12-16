package action

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"unsafe"

	"nvim-help/internal/utils"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/pkg/archive"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/go-connections/nat"
	"github.com/moby/moby/client"
)

const (
	start = "start"
	stop  = "stop"

	imageName         = "debug/go"
	argsGoVersionName = "GO_VERSION"
)

var (
	_ Action = (*debugByDocker)(nil)

	ErrIncorrectRequest = errors.New("action is not valid")

	delveMount = mount.Mount{
		Type:   mount.TypeBind,
		Source: os.Getenv("GOPATH") + "/bin/delve",
		Target: "/root/delve",
	}
)

type debugByDocker struct {
	cli  client.APIClient
	fs   *flag.FlagSet
	args *string
}

func NewDebugByDocker() Action {
	dd := &debugByDocker{fs: flag.NewFlagSet("debug-docker", flag.ContinueOnError)}
	dd.args = dd.fs.String("json", "{}", "args in json format.")
	return dd
}

// Action implements Action.
func (d *debugByDocker) Action() string {
	return "debug-docker"
}

// Run implements Action.
func (d *debugByDocker) Run(args []string) *Result {
	if len(args) < 1 {
		return NewFailResult(ErrIncorrectRequest)
	}

	var err error
	if d.cli, err = client.NewClientWithOpts(client.FromEnv); err != nil {
		return NewFailResult(err)
	}

	var (
		ctx = context.Background()
		act = args[0]
	)

	if len(args) > 1 {
		if err = d.fs.Parse(args[1:]); err != nil {
			return NewFailResult(err)
		}
	}

	switch act {
	case "start":
		resp, err := d.Start(ctx, *d.args)
		if err != nil {
			return NewFailResult(err)
		}
		return NewSuccessResult(resp)
	case "stop":
		if err = d.Stop(ctx, *d.args); err != nil {
			return NewFailResult(err)
		}
		return NewSuccessResult(nil)
	case "build":
		return nil
	default:
		return NewFailResult(ErrIncorrectRequest)
	}
}

// Usege implements Action.
func (d *debugByDocker) Usege() string {
	panic("unimplemented")
}

type startRequset struct {
	ProjectPath string `json:"project_path"`
}

type startResponse struct {
	Port uint16 `json:"port"`
}

func (d *debugByDocker) Start(ctx context.Context, request string) (resp *startResponse, err error) {
	req := &startRequset{}
	if err = json.Unmarshal(unsafe.Slice(unsafe.StringData(request), len(request)), req); err != nil {
		return nil, err
	}

	projectName := path.Base(req.ProjectPath)
	if err = d.deleteContainer(ctx, projectName); err != nil {
		return nil, err
	}

	containList, err := d.cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("ancestor", imageName)),
	})
	if err != nil {
		return nil, err
	}

	port := uint16(38697)
	pubPort := make(map[uint16]struct{})
	for _, c := range containList {
		for _, port := range c.Ports {
			pubPort[port.PublicPort] = struct{}{}
		}
	}
	for {
		if _, ok := pubPort[port]; !ok {
			break
		}
		port++
	}

	ports, portBinds, err := nat.ParsePortSpecs([]string{fmt.Sprintf("%d:%d", port, port)})
	if err != nil {
		return nil, err
	}

	workPath := path.Join("/root", path.Base(req.ProjectPath))
	mounts := make([]mount.Mount, 0, 2)
	mounts = append(mounts, delveMount)
	mounts = append(mounts, mount.Mount{
		Type:   mount.TypeBind,
		Source: utils.ConvertEnvPlace(req.ProjectPath),
		Target: workPath,
	})

	respCC, err := d.cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        imageName,
			ExposedPorts: ports,
			WorkingDir:   workPath,
		},
		&container.HostConfig{
			Mounts:       mounts,
			PortBindings: portBinds,
		},
		nil, nil,
		projectName,
	)
	if err != nil {
		return nil, err
	}

	if err := d.cli.ContainerStart(ctx, respCC.ID, container.StartOptions{}); err != nil {
		return nil, err
	}

	return &startResponse{
		Port: port,
	}, nil
}

type stopRequest struct {
	ProjectPath string `json:"project_path"`
}

func (d *debugByDocker) Stop(ctx context.Context, request string) (err error) {
	req := &stopRequest{}
	if err = json.Unmarshal(unsafe.Slice(unsafe.StringData(request), len(request)), req); err != nil {
		return err
	}

	if err = d.deleteContainer(ctx, path.Base(req.ProjectPath)); err != nil && !strings.Contains(
		err.Error(),
		"Is the docker daemon running",
	) {
		return err
	}

	return nil
}

type buildRequest struct {
	ProjectPath    string `json:"project_path"`
	DockerfilePath string `json:"dockerfile_path"`
}

func (d *debugByDocker) Build(ctx context.Context, request string) (err error) {
	req := &buildRequest{}
	if err = json.Unmarshal(unsafe.Slice(unsafe.StringData(request), len(request)), req); err != nil {
		return err
	}

	goVersion, err := d.getGoVersionByProject(req.ProjectPath)
	if err != nil {
		return err
	}

	return d.build(ctx, req.DockerfilePath, goVersion)
}

func (d *debugByDocker) deleteContainer(ctx context.Context, containerName string) (err error) {
	containers, err := d.cli.ContainerList(
		ctx,
		container.ListOptions{
			All:     true,
			Filters: filters.NewArgs(filters.Arg("name", containerName)),
		},
	)
	if err != nil {
		return err
	}

	var errs error
	for _, c := range containers {
		if err = d.cli.ContainerRemove(ctx, c.ID, container.RemoveOptions{Force: true}); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return nil
}

func (d *debugByDocker) build(ctx context.Context, dockerfilePath string, goVersion string) (err error) {
	reader, err := archive.TarWithOptions(dockerfilePath, &archive.TarOptions{})
	if err != nil {
		return err
	}

	args := make(map[string]*string, 0)
	args[argsGoVersionName] = &goVersion
	resp, err := d.cli.ImageBuild(ctx, reader, types.ImageBuildOptions{
		Tags:      []string{imageName + goVersion},
		BuildArgs: args,
	})
	if err != nil {
		return err
	}

	// TODO stream
	_, err = io.ReadAll(resp.Body)

	return err
}

func (d *debugByDocker) getGoVersionByProject(ProjectPath string) (goVersion string, err error) {
	modPath, err := utils.GetModPath(ProjectPath)
	if err != nil {
		return "", err
	}

	return utils.GetGoVersion(path.Join(modPath, "go.mod"))
}
