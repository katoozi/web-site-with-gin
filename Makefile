# note: call scripts from /scripts
all: help
help:
	@echo "Usage:"
	@echo "    make build                 build docker stack file"
	@echo "    make up                    run the services"
	@echo "    make stop                  stop the services"
	@echo "    make deploy                deploy the stack   => make stack_name=postgres stack_file=postgres/docker-stack.yml deploy"
	@echo "    make rebuild               rebuild service(s) => make services='postgres web-server redis' rebuild"

build:
	docker-compose --compatibility -f docker-compose.yml config > docker-stack.yml

up:
	docker-compose -f docker-compose.yml --compatibility up --build -d

rebuild:
	docker-compose -f docker-compose.yml --compatibility up -d --no-deps --build ${services}

down:
	docker-compose -f docker-compose.yml --compatibility down

deploy:
	docker stack deploy -c ${stack_file} ${stack_name}
