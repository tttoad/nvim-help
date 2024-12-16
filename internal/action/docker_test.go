package action

import (
	"context"
	"fmt"
	"testing"
)

func TestBuild(t *testing.T) {
	d := debugByDocker{}
	d.Run([]string{"build"})
	fmt.Println(d.build(context.TODO(), "/Users/toad/work/nvim-help/internal", "/Users/toad/work/nvim-help/Dockerfile"))
}
