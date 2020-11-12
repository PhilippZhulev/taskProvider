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
	go build -o ./build/services/broker/bin/broker -v ./services/broker/cmd/broker
	# copy broker config
	cp ./services/broker/configs/broker.toml ./build/services/broker
	# getway
	go build -o ./build/services/getway/bin/getway -v ./services/getway/cmd/getway
	# copy getway config
	#cp ./services/getway/configs/getway.toml ./build/services/getway/configs
	# user
	go build -o ./build/services/user/bin/user -v ./services/user/cmd/user
	# copy user config
	cp ./services/user/configs/user.toml ./build/services/user
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

.PHONY: build_d
build_d:
	docker build -t task-provider-user ./services/user   
	docker build -t task-provider-getway ./services/getway 
	docker build -t task-provider-broker ./services/broker 

.PHONY: start_d
start_d:
	docker run -p 8081:8081 -it --rm --name getway task-provider-getway
	docker run -p 5040:5040 -it --rm --name broker task-provider-broker
	docker run -p 5041:5041 -it --rm --name user task-provider-user


.DEFAULT_GOAL := protogen