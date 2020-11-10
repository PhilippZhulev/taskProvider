package protogenerator

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"sync"

	"gitlab.com/taskProvider/protogen/internal/app/configure"
)
var wg sync.WaitGroup

// NewProto ...
func NewProto(conf *configure.Conf) {
	for _, despItem := range conf.Desp {
		wg.Add(1)
		go outIteration(conf, despItem)
	}

	wg.Wait()
}

func outIteration(conf *configure.Conf, item configure.Desp) {
	for _, outItem := range item.Out {

		cmd := exec.Command(
			"protoc", 
			"-I.", 
			"-I" + conf.Dir,
			"--go_out=" + outItem,
			"--go_opt=paths=source_relative",
			"--go-grpc_out=" + outItem,
			"--go-grpc_opt=paths=source_relative",
			item.In,
		)

		if item.Getway {
			cmd = exec.Command(
				"protoc", 
				"-I.", 
				"-I" + conf.Dir,
				"--go_out=" + outItem,
				"--go_opt=paths=source_relative",
				"--grpc-gateway_opt",
				"paths=source_relative", 
				"--grpc-gateway_out=logtostderr=true:" + outItem,
				"--go-grpc_out=" + outItem,
				"--go-grpc_opt=paths=source_relative",
				item.In,
			)
		}

		stderr, _ := cmd.StderrPipe()
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		go func(s io.ReadCloser) {
			scanner := bufio.NewScanner(s)
			for scanner.Scan() {
				m := scanner.Text()
				log.Println(m)
			}
			wg.Done()
		}(stderr)
	}
}