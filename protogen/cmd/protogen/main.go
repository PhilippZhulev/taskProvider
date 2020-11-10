package main

import (
	"gitlab.com/taskProvider/protogen/internal/app/configure"
	"gitlab.com/taskProvider/protogen/internal/app/protogenerator"
)

func main() {
	conf := configure.NewConf()
	protogenerator.NewProto(conf);
}