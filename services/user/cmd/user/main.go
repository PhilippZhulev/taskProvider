package main

import (
	"flag"
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

	wordPtr := flag.String("conf", dir + "/../configs/user.toml", "config url")
	
	flag.Parse()

	config := configurator.NewConfig()
	_, err = toml.DecodeFile(*wordPtr, config)
	if err != nil {
		log.Fatalf("failed: %v", err)
	}

	err = server.Run(config)
	if err != nil {
		log.Fatalf("failed: %v", err)
	}
}