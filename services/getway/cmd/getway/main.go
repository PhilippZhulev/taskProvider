package main

import (
	"flag"
	"log"

	"gitlab.com/taskProvider/services/getway/internal/app/server"

	"github.com/golang/glog"
)

const (
	address     = "localhost:5050"
	defaultName = "world"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	  
	if err := server.NewServer(); err != nil {
		log.Fatal(err)
	}
}