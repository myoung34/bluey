GOOS = darwin
GOARCH = amd64

test:
	go test -coverprofile coverage.out ./...
	go tool cover -func=coverage.out


build:
	GOOS=linux GOARCH=arm GODEBUG=cgocheck=0 go build -o dist/bluey
	GOOS="" GOARCH="" GODEBUG=cgocheck=0 go build -o dist/bluey

run:
	sudo ./dist/bluey -c test.ini

build-docker:
	docker-compose build

run-docker: build-docker
	docker-compose run bluey
