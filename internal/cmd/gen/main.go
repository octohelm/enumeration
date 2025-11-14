package main

import (
	"context"

	"github.com/octohelm/gengo/pkg/gengo"
	"github.com/octohelm/x/logr"
	"github.com/octohelm/x/logr/slog"
)

import (
	_ "github.com/octohelm/enumeration/devpkg/enumgen"
)

func main() {
	c, err := gengo.NewContext(&gengo.GeneratorArgs{
		Entrypoint: []string{
			"github.com/octohelm/enumeration/testdata/model",
		},
		OutputFileBaseName: "zz_generated",
		Globals: map[string][]string{
			"gengo:runtimedoc": {},
		},
	})
	if err != nil {
		panic(err)
	}

	ctx := logr.WithLogger(context.Background(), slog.Logger(slog.Default()))
	if err := c.Execute(ctx, gengo.GetRegisteredGenerators()...); err != nil {
		panic(err)
	}
}
