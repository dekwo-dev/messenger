# syntax=docker/dockerfile:1

FROM golang:alpine3.19
WORKDIR /app
COPY . .
CMD ["go", "run", "."]
