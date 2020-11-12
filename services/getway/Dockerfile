FROM golang:1.15

WORKDIR /go/src/taskProvider/services/getway
COPY . .

WORKDIR /go/src/taskProvider
RUN go mod init gitlab.com/taskProvider
RUN go mod tidy
RUN go get -d -v ./...
RUN go install -v ./...

CMD go run ./services/getway/cmd/getway

EXPOSE 8081