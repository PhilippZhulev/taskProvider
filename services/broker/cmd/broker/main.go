package main

import (
	"log"
	"os"
	"path/filepath"

	"gitlab.com/taskProvider/services/broker/internal/app/server"

	"gitlab.com/taskProvider/services/broker/internal/app/configurator"

	"github.com/BurntSushi/toml"
)


func main() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	
	config := configurator.NewConfig()
	_, err = toml.DecodeFile(dir + "/../configs/broker.toml", config)
	if err != nil {
		log.Fatal(err)
	}

	logs, err := server.NewServer(config)
	if err != nil {
		logs.Fatalf("failed: %v", err)
	}
}