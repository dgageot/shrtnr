# sudo ifconfig lo0 alias 127.0.0.2 up

run:
	docker build -t dgageot/shrtnr .
	-docker rm -f shrtnr
	docker run --rm --name shrtnr -d -p 127.0.0.2:80:8080 -v $(HOME)/links:/links dgageot/shrtnr