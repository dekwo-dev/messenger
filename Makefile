build:
	docker build -t guestbook-backend .
dev:
	docker run -dp 127.0.0.1:8000:8000 guestbook-backend
up:
	docker compose up
down:
	docker compose down
