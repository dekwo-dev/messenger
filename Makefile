proto:
	protoc --go_out=./ --proto_path=./ ./def/event.proto ./def/log.proto
build:
	docker build -t dekwo.dev/messenger .
up:
	docker compose up
down:
	docker compose down
clean:
	docker image prune -a
fmt:
	golines -w *.go
