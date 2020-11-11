package main

import (
	"flag"
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
		log.Fatalf("failed: %v", err)
	}

	wordPtr := flag.String("conf", dir + "/../configs/broker.toml", "config url")
	
	flag.Parse()
	
	config := configurator.NewConfig()
	_, err = toml.DecodeFile(wordPtr, config)
	if err != nil {
		log.Fatal(err)
	}

	logs, err := server.NewServer(config)
	if err != nil {
		logs.Fatalf("failed: %v", err)
	}
}