TARGET = binny
SOURCES = statik/statik.go $(wildcard *.go)

.PHONY: build
build: $(TARGET)

.PHONY: run
run: $(TARGET)
	./binny echo hello

.PHONY: clean
clean:
	docker rmi binny --force 2>/dev/null
	rm -rf binny assets/ statik/

$(TARGET): $(SOURCES)
	go build -o binny

statik/statik.go: assets/image.tar.gz
	statik -src=$(abspath assets)

assets/image.tar.gz: Dockerfile
	docker build -t binny .
	mkdir -p assets
	docker save binny | gzip > assets/image.tar.gz
	docker rmi binny
