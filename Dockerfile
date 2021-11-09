FROM golang:1.17.3-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

# see dockerignore
COPY . .

RUN go build -o /app/main

CMD [ "/app/main" ]