TARGET = binny
SOURCES = Dockerfile image.tar.gz pkged.go $(wildcard *.go)

.PHONY: build
build: $(TARGET)

.PHONY: run
run: $(TARGET)
	./binny
	docker run --rm binny hello
	docker rmi binny

.PHONY: clean
clean:
	rm -f image.tar.gz pkged.go $(TARGET)

$(TARGET): $(SOURCES)
	go build -o binny main.go

image.tar.gz: Dockerfile
	docker build -t binny .
	docker save binny | gzip > image.tar.gz
	docker rmi binny

pkged.go: image.tar.gz
	pkger
