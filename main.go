package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/HomayoonAlimohammadi/mini-grpc-gateway/config"
	"github.com/HomayoonAlimohammadi/mini-grpc-gateway/export"
)

func main() {
	var flags flag.FlagSet

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		conf, err := config.ExtractServiceConfig(gen, config.ServicesConfigJSON)
		if err != nil {
			panic(err)
		}

		export.Export(gen, conf)

		return nil
	})
}
