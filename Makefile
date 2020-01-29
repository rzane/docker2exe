.PHONY: binny
build: binny

.PHONY: run
run: pkged.go
	go run main.go

image.tar.gz: Dockerfile
	docker build -t binny .
	docker save binny | gzip > image.tar.gz
	docker rmi binny

pkged.go: image.tar.gz
	pkger

binny: pkged.go *.go
	go build -o binny main.go
