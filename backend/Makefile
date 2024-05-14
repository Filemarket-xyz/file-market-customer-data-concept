include .env

run:
	go run cmd/app/main.go -cfg configs/local

postgres.start:
	if [ ! "$(shell docker ps -q -f name=fmconcept-postgres)" ]; then \
        if [ "$(shell docker ps -aq -f status=exited -f name=fmconcept-postgres)" ]; then \
            docker rm fmconcept-postgres; \
        fi; \
		docker run --restart unless-stopped -d -p 1339:5432 -e POSTGRES_PASSWORD=1337 --name fmconcept-postgres postgres:16; \
        sleep 5; \
    fi;
	-docker exec fmconcept-postgres psql -U postgres -c "create user fmconcept password '1337'"
	-docker exec fmconcept-postgres psql -U postgres -c "create database fmconcept"
	-docker exec fmconcept-postgres psql -U postgres -c "grant all on database fmconcept to fmconcept"
	-docker exec fmconcept-postgres psql -U postgres -c "ALTER DATABASE fmconcept OWNER TO fmconcept"
	-docker exec fmconcept-postgres psql -U postgres -c "GRANT ALL ON SCHEMA PUBLIC TO fmconcept"

postgres.stop:
	docker stop fmconcept-postgres
	docker rm fmconcept-postgres

redis.start:
	if [ ! "$(shell docker ps -q -f name=fmconcept-redis)" ]; then \
   	 	docker run --name fmconcept-redis -d -p $(REDIS_PORT):6379 redis; \
    fi; \
	docker start fmconcept-redis

redis.flush:
	docker exec fmconcept-redis redis-cli FLUSHALL

redis.stop:
	docker stop fmconcept-redis
	docker rm fmconcept-redis

build-models:
	docker run --rm --user $(shell id -u):$(shell id -g) -e GOPATH=$(go env GOPATH):/go -v ${HOME}:${HOME} -w $(shell pwd) quay.io/goswagger/swagger generate model --spec=../swagger.yaml
	go mod tidy

migrate.up:
	goose -dir migrations postgres "host=localhost port=1339 user=fmconcept password=1337 dbname=fmconcept sslmode=disable" up

migrate.down:
	goose -dir migrations postgres "host=localhost port=1339 user=fmconcept password=1337 dbname=fmconcept sslmode=disable" down

migrate.status:
	goose -dir migrations postgres "host=localhost port=1339 user=fmconcept password=1337 dbname=fmconcept sslmode=disable" status

tests:
	go test -v ./integration_tests/...