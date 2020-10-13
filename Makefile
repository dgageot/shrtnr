IMG ?= dgageot/shrtnr
NAME ?= shrtnr
LINKS_HOME ?= $(HOME)/links

all: build run

.PHONY: build
build:
	docker build -t $(IMG) .
	
.PHONY: run
run: setuplo0
	-docker rm -f $(NAME)
	docker run --rm --name $(NAME) -d -p 127.0.0.2:80:8080 -v $(LINKS_HOME):/root/links $(IMG)

.PHONY: run-local
run-local:
	# Try http://localhost:8080/google in a browser.
	go run main.go

.PHONY: setuplo0
setuplo0:
	# Make sure 127.0.0.2 is configured as a link-local address.
	@ifconfig lo0 | grep -q "127.0.0.2" || sudo ifconfig lo0 alias 127.0.0.2 up
