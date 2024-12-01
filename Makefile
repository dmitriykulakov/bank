.PHONY: run_server

all: run_server

run_server: 
	sudo docker-compose -f ./docker-compose.yaml up --build 
	
test:
	go test ./ -v -count=1
	
	