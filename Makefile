TARGET = binny
SOURCES = $(wildcard *.go)

.PHONY: build
build: $(TARGET)

.PHONY: run
run: $(TARGET)
	./binny echo hello

.PHONY: image
image:
	mkdir -p assets
	docker build -t binny .
	docker save binny | gzip > assets/image.tar.gz
	docker rmi binny

$(TARGET): $(SOURCES)
	go build -o binny
