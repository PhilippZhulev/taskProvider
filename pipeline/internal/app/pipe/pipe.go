package pipe

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"gitlab.com/taskProvider/pipeline/internal/app/bot"
	"gitlab.com/taskProvider/pipeline/internal/app/configure"
)

// NewPipeline ...
func NewPipeline(conf *configure.Conf, command, errChan chan string, init *bot.Init) func() {

	// project info
	fmt.Println("\033[1;36m", conf.Name, "services", "----->", "version:", conf.Version, "\033[0m")

	// read path to services
	var cmd []*exec.Cmd

	cmd = pipe(conf, cmd, errChan)
	init.State = true

	// wait and command
    for c := range command {
        // command to exit program
        if c == "q" {
			offPipeline(cmd)
			init.State = false
			log.Println("Pipeline off.")
		}
		
		if c == "s" {
			var new []*exec.Cmd
			cmd = pipe(conf, new, errChan)
			init.State = true
			log.Println("Pipeline on.")
		}
		if c == "r" {
			offPipeline(cmd)
			var new []*exec.Cmd
			cmd = pipe(conf, new, errChan)
			log.Println("Pipeline restart.")
		}
	}

	// is main fun and, then kill active process
	return func() {
		offPipeline(cmd)
	}
}

func pipe(conf *configure.Conf, cmd []*exec.Cmd, errChan chan string) []*exec.Cmd {
	// init log name creater
	var lst = replLog(conf.Services)

	for i, item := range conf.Services {
		// run service
		cmd = append(cmd, exec.Command(item[1]))
		stdout, _ := cmd[i].StdoutPipe()
		stderr, _ := cmd[i].StderrPipe()
		if err := cmd[i].Start(); err != nil {
			log.Fatal(err) 
		}

		// scan log output
		go func(s io.ReadCloser, name string) {
			scanner := bufio.NewScanner(s)
			for scanner.Scan() {
				m := scanner.Text()
				fmt.Println("\033[1;33mLOG", lst(name), ":\033[0m", m)

				// check is Log true then write errors by file
				if conf.ErrorLog {
					writeErrFile(conf.LogURL, "LOG " + lst(name) + ": " + m)
				}
			}
		}(stdout, item[0])

		// scan err output
		go func(e io.ReadCloser, name string) {
			scanner := bufio.NewScanner(e)
			for scanner.Scan() {
				m := scanner.Text()
				fmt.Println("\033[1;31mERROR", lst(name), ":\033[0m", m)
				errChan <- "ERROR " + lst(name) + ": " + m
				// check is ErrorLog true then write errors by file
				if conf.ErrorLog {
					writeErrFile(conf.ErrorLogURL, "ERROR " + lst(name) + ": " + m)
				}
			}
		}(stderr, item[0])
	}

	return cmd
}

//off
func offPipeline(cmd []*exec.Cmd) {
	for _, c := range cmd {
		cmd := exec.Command("kill", strconv.Itoa(c.Process.Pid))
		if err := cmd.Start(); err != nil {
			log.Fatal(err) 
		}
	}
}

// write log file
func writeErrFile(f, text string) {
    file, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
     
    if err != nil{
        fmt.Println("Unable to create log file:", err) 
        os.Exit(1) 
    }
    defer file.Close() 
    file.WriteString(text + "\n")
}

// create log name
func replLog(e [][]string) func(string) string {
	var max int 
	var spaceArr []string

	for ind := 0; ind < (len(e) - 1); ind++ {
		if len(e[ind][0]) > len(e[ind + 1][0]) {
			max = len(e[ind][0])
		}
	}

	return func (s string) string  {
		for i := 0; i < max; i++ {
			spaceArr = append(spaceArr, " ")
		}

		var result string
		for i := 0; i < len(spaceArr); i++ {
			if i < len(s) {
				spaceArr[i] = string(s[i])
			}
			result += spaceArr[i]
		}
		spaceArr = spaceArr[:0]
		return result
	}
}