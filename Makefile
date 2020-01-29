TARGET = binny
SOURCES = Dockerfile image.tar.gz pkged.go $(wildcard *.go)

.PHONY: build
build: $(TARGET)

.PHONY: run
run: $(TARGET)
	./binny hello

.PHONY: clean
clean:
	docker rmi -f binny
	rm -f image.tar.gz pkged.go $(TARGET)

$(TARGET): $(SOURCES)
	go build -o binny main.go

image.tar.gz: Dockerfile
	docker build -t binny .
	docker save binny | gzip > image.tar.gz
	docker rmi binny

pkged.go: image.tar.gz
	pkger
