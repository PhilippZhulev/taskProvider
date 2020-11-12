.PHONY: rdsb
rdsb:
	docker-compose up --build  

.PHONY: rds
rds:
	docker-compose up 

.PHONY: rdsstop
rdsstop:
 	docker-compose down 

.PHONY: protogen
protogen:
	./protogen/bin/protogen

.PHONY: build
build:
	go build -o ./services/broker/bin -v ./services/broker/cmd/broker
	go build -o ./services/getway/bin -v ./services/getway/cmd/getway
	go build -o ./services/user/bin -v ./services/user/cmd/user


.PHONY: build_p
bup:
	go build -o ./build/services/broker/bin/broker -v ./services/broker/cmd/broker
	cp ./services/broker/configs/broker.toml ./build/services/broker
	go build -o ./build/services/getway/bin/getway -v ./services/getway/cmd/getway
	#cp ./services/getway/configs/getway.toml ./build/services/getway/configs
	go build -o ./build/services/user/bin/user -v ./services/user/cmd/user
	cp ./services/user/configs/user.toml ./build/services/user
	go build -o ./build/pipeline/pipeline -v ./pipeline/cmd/pipeline
	cp ./pipeline/pipeconf.json ./build/pipeline

.PHONY: buildpipe
buildpipe:
	go build -o ./pipeline/bin -v ./pipeline/cmd//pipeline

.PHONY: pipeline
pipeline:
	go run ./pipeline/cmd/pipeline    

.PHONY: pgb
pgb:
	go build -o ./protogen/bin -v ./protogen/cmd/protogen


.DEFAULT_GOAL := protogen