FROM golang:1.18

WORKDIR /app

COPY Makefile ./Makefile
COPY cmd ./cmd
COPY pkg ./pkg
COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN make build
