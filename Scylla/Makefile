# Makefile

install: go-mod db-prepare

go-mod:
	go mod tidy

db-prepare:
	rm -rf data
	docker compose up -d

start:
	docker exec -it cdb01 sh -c 'cockroach init --insecure --host=cdb01:26357'

stop:
	docker compose down
	rm -rf data

secondary_indexes:
	go run main.go --script-name=secondary_indexes --is-show-all=false --is-need-clear=true --secret-key="my secret key"

view:
	go run main.go --script-name=views --view-name=researchView --is-show-all=false --is-need-clear-table=true
