## taskProvider


#### Make commands

```
// stop docker services
make rdsstop

// run services to docker and build
make rdsb

// run services to docker
make rds

// run services pipeline
make pipeline

// build app
make bup

// build pipeline component
make buildpipe

// build services
make build

// generate proto files
make protogen

// build protogen component
make pgb
```

#### Commands

```
// path
export PATH="$PATH:$(go env GOPATH)/bin"

// change os linux
export GOOS=linux  
export GOARCH=amd64

// change os mac
export GOOS=darwin  
export GOARCH=amd64

// get containers ip
docker inspect --format='{{.Name}} - {{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(docker ps -aq)

```