FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV HTTP_USER_PASSWORD=rndPassword123
ENV SECRET_JWT_KEY=hugYourNewJWTKey

RUN go build -o app ./cmd/app/main.go

EXPOSE 8080

CMD ["./main"]

