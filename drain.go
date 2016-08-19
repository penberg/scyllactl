package main

import (
	"fmt"
	"github.com/mkideal/cli"
	"net/http"
)

var _ = app.Register(&cli.Command{
	Name: "drain",
	Desc: "Drain the node (stop accepting writes and flush all column families)",
	Argv: func() interface{} { return new(drainT) },
	Fn:   drain,
})

type drainT struct {
	cli.Helper
	Host string `cli:"host" usage:"Node hostname or IP address" dft:"localhost"`
	Port string `cli:"port" usage:"API server port number" dft:"10000"`
}

func drain(ctx *cli.Context) error {
	argv := ctx.Argv().(*drainT)
	baseURL := fmt.Sprintf("http://%s:%s", argv.Host, argv.Port)
	if _, err := http.Post(baseURL + "/storage_service/drain", "application/json", nil); err != nil {
		return err
	}
	return nil
}
