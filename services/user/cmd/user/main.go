package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"gitlab.com/taskProvider/services/user/internal/app/configurator"
	"gitlab.com/taskProvider/services/user/internal/app/server"
)

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("failed: %v", err)
	}
	
	config := configurator.NewConfig()
	_, err = toml.DecodeFile(dir + "/../configs/user.toml", config)
	if err != nil {
		log.Fatalf("failed: %v", err)
	}

	err = server.Run(config)
	if err != nil {
		log.Fatalf("failed: %v", err)
	}
}