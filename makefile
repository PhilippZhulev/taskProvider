.PHONY: protoinit
protoinit:
	export PATH="$PATH:$(go env GOPATH)/bin"


.PHONY: protogen
protogen:

	./protogen/bin/protogen

.PHONY: build
build:
	# broker
	go build -o ./services/broker/bin -v ./services/broker/cmd/broker
	# getway
	go build -o ./services/getway/bin -v ./services/getway/cmd/getway
	# user
	go build -o ./services/user/bin -v ./services/user/cmd/user


.PHONY: build_p
build_p:
	# broker
	go build -o ./build/services/broker/bin -v ./services/broker/cmd/broker
	# copy broker config
	cp ./services/broker/configs/broker.toml ./build/services/broker/configs
	# getway
	go build -o ./build/services/getway/bin -v ./services/getway/cmd/getway
	# copy getway config
	#cp ./services/getway/configs/getway.toml ./build/services/getway/configs
	# user
	go build -o ./build/services/user/bin -v ./services/user/cmd/user
	# copy user config
	cp ./services/user/configs/user.toml ./build/services/user/configs
	# pipeline
	go build -o ./build/pipeline/pipeline -v ./pipeline/cmd/pipeline
	# copy pipline config
	cp ./pipeline/pipeconf.json ./build/pipeline
	

.PHONY: serve
serve:

	# broker
	go build -o ./services/broker/bin -v ./services/broker/cmd/broker
	# getway
	go build -o ./services/getway/bin -v ./services/getway/cmd/getway
	# user
	go build -o ./services/user/bin -v ./services/user/cmd/user

	# run getway
	./services/getway/bin/getway
	# run broker
	./services/broker/bin/broker
	# run user
	./services/user/bin/user


.PHONY: start
start:

	# run getway
	./services/getway/bin/getway
	# run broker
	./services/broker/bin/broker
	# run user
	./services/user/bin/user


.PHONY: broker
broker:
	# run broker
	./services/broker/bin/broker


.PHONY: getway
getway:
	# run getway
	./services/getway/bin/getway


.PHONY: user
user:
	# run user
	./services/user/bin/user


.PHONY: buildpipe
buildpipe:
	# build pipeline
	go build -o ./pipeline/bin -v ./pipeline/cmd//pipeline


.PHONY: pipeline
pipeline:
	# run pipeline
	go run ./pipeline/cmd/pipeline    

.PHONY: buildpg
buildpg:
	# build pipeline
	go build -o ./protogen/bin -v ./protogen/cmd/protogen


.PHONY: linux
linux:
	export GOOS=linux  
	export GOARCH=amd64


.PHONY: mac
mac:
	export GOOS=darwin  
	export GOARCH=amd64

.PHONY: d_build
d_build:
	docker build -t task-provider .

.PHONY: d_start
d_start:
	docker run -it --rm --name my-running-app task-provider


.DEFAULT_GOAL := protogen