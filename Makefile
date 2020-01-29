.PHONY: binny
build: binny

.PHONY: run
run: pkged.go
	go run main.go

image.tar: Dockerfile
	docker build -t binny .
	docker save -o image.tar binny
	docker rmi binny

pkged.go: image.tar
	pkger

binny: pkged.go *.go
	go build -o binny main.go
