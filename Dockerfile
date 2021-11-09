FROM golang:1.17.3-alpine

# needed for go modules
ENV GO111MODULE=on

# sets env vars for host
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

# gets go modules
RUN go mod tidy -v
RUN go mod download

# see dockerignore
COPY . .

RUN GOARCH=$TARGETARCH GOOS=$TARGETOS go build -o /app/main

CMD [ "/app/main" ]