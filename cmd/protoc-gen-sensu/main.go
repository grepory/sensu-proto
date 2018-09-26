package main

import (
	"github.com/gogo/protobuf/vanity/command"
	"github.com/sensu/sensu-proto/plugin/defaulter"
)

func main() {
	// run defaulter
	resp := command.GeneratePlugin(
		command.Read(),
		defaulter.NewDefaulter(),
		"_defaults.pb.go",
	)
	command.Write(resp)
}
