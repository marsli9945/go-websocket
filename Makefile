build:
	docker build -t go-websocket .
rebuild:
	docker stop go-websocket
	docker build -t go-websocket .
	docker rm go-websocket
	docker run -dit -p 7777:7777 --name=go-websocket go-websocket
run:
	docker run -dit -p 7777:7777 --name=go-websocket go-websocket
start:
	docker start go-websocket
restart:
	docker restart go-websocket
stop:
	docker stop go-websocket