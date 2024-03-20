proto:
	protoc --go_out=./ --proto_path=./ ./def/event.proto ./def/log.proto
build:
	docker build -t dekwo.dev/messenager .
up:
	docker compose up
down:
	docker compose down
fmt:
	golines -w *.go
