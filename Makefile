IMAGE_NAME = ggp

.PHONY: all build clean install

all: build clean

build:
	mkdir -p binary
	go build -a  \
		-gcflags=all="-l -B" \
		-ldflags="-w -s" \
		-o binary/$(IMAGE_NAME) \
		./...

install: build
	sudo cp binary/$(IMAGE_NAME) /usr/local/bin

clean:
	rm -rf binary
