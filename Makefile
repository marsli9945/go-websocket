build:
	docker build -t go-websocket .
rebuild:
	docker stop go-websocket
	docker rm go-websocket
	docker build -t go-websocket .
	docker run --env GAPI_HOST='https://gapics.touch4.me' --env GAPI_CLIENT_ID='net5ijy' --env GAPI_CLIENT_SECRET='y4cZe@wmGBofIlQ' --env GAPI_USERNAME='lilei@tuyoogame.com' --env GAPI_PASSWORD='Lilei@123' -dit -p 7777:7777 --name=go-websocket go-websocket
	docker start go-websocket
run:
	docker run --env GAPI_HOST='https://gapics.touch4.me' --env GAPI_CLIENT_ID='net5ijy' --env GAPI_CLIENT_SECRET='y4cZe@wmGBofIlQ' --env GAPI_USERNAME='lilei@tuyoogame.com' --env GAPI_PASSWORD='Lilei@123' -dit -p 7777:7777 --name=go-websocket go-websocket
start:
	docker start go-websocket
restart:
	docker restart go-websocket
stop:
	docker stop go-websocket
log:
	docker logs -f go-websocket