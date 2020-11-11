FROM golang:1.15

WORKDIR /go/src/taskProvider
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 5000

CMD go run ./services/user/cmd/user -conf=./services/user/configs/user.toml