package main

import (
	"fmt"
	"log"

	"github.com/keploy/go-sdk/integrations/kfasthttp"
	"github.com/keploy/go-sdk/keploy"
	"github.com/valyala/fasthttp"
)

const port = "3639"

var kConfig *keploy.Keploy

func init() {
	kConfig = keploy.New(keploy.Config{
		App: keploy.AppConfig{
			Name: "gfstat",
			Port: port,
		},
		Server: keploy.ServerConfig{}, // default
	})
}

func serve() {
	kMiddleware := kfasthttp.FastHttpMiddleware(kConfig)
	routes := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/":
		default:
			ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		}
	}

	fmt.Println("Listening on port: ", port)
	log.Fatal(fasthttp.ListenAndServe(":"+port, kMiddleware(routes)))
}
