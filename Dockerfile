FROM golang:1.19

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY api ./api
COPY config ./config
COPY jwt ./jwt
COPY logger ./logger
COPY mail ./mail
COPY store ./store
COPY .env .gitignore ./
COPY .golangci-lint.yaml docker-compose.yaml ./
COPY Dockerfile main.go ./

RUN go build -o /awesome-api

EXPOSE 3000

CMD [ "/awesome-api" ]