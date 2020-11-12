FROM golang:1.15

WORKDIR /go/src/taskProvider/services/user
COPY . .

WORKDIR /go/src/taskProvider
RUN go mod init gitlab.com/taskProvider
RUN go mod tidy
RUN go get -d -v ./...
RUN go install -v ./...

CMD go run ./services/user/cmd/user -conf=./services/user/configs/user.toml

EXPOSE 5041