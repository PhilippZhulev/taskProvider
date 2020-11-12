FROM golang:1.15

WORKDIR /go/src/taskProvider/services/broker
COPY . .

WORKDIR /go/src/taskProvider
RUN go mod init gitlab.com/taskProvider
RUN go mod tidy
RUN go get -d -v ./...
RUN go install -v ./...

CMD go run ./services/broker/cmd/broker -conf=./services/broker/configs/broker.toml

EXPOSE 7040