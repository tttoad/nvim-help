package action

import (
	"context"
	"fmt"
	"testing"

	"github.com/moby/moby/client"
)

func TestChange(t *testing.T) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	fmt.Println(cli, err)
	err = cli.ContainerKill(context.Background(), "", "")
	fmt.Println(err)
}
