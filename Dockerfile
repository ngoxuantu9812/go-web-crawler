FROM golang:1.20-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go


EXPOSE 8080
CMD [ "/app/main" ]
