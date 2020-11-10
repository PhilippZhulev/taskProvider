package main

import (
	"fmt"

	"gitlab.com/taskProvider/pipeline/internal/app/bot"
	"gitlab.com/taskProvider/pipeline/internal/app/configure"
	"gitlab.com/taskProvider/pipeline/internal/app/pipe"
)

func main() {
	init := bot.Init{}
	
	command := make(chan string) // command chanel
	errors := make(chan string) // errors chanel

	// decode json to struct
	conf := configure.NewConf()

	// telegram bot
	go bot.NewBot(conf, command, errors, &init)

	// pipeline
	close := pipe.NewPipeline(conf, command, errors, &init)
	
	var input string
	fmt.Scanln(&input)
	
	defer close()
}
