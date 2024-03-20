proto:
	protoc --go_out=./ --proto_path=./ ./def/messager.proto ./def/log.proto
tuild:
	docker build -t guestbook-backend .
dev:
	docker run -dp 127.0.0.1:8000:8000 guestbook-backend
up:
	docker compose up
down:
	docker compose down
fmt:
	golines -w *.go
