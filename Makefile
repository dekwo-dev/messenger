proto:
	protoc --go_out=./ --proto_path=./ ./def/messager.proto ./def/log.proto
build:
	docker build -t guestbook-backend .
up:
	docker compose up
down:
	docker compose down
fmt:
	golines -w *.go
