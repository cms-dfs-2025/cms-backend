FROM golang:1.24.2-bookworm

WORKDIR /usr/src/app

COPY go.mod go.sum .
RUN go mod download

COPY *.go .
RUN go build -v -o /usr/local/bin/app .

CMD ["app"]
