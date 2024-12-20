FROM golang:1.23.3-alpine

RUN apk add --no-cache \
    build-base \
    sqlite-dev

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./
RUN go build -o .

EXPOSE 8080

CMD ["/events-api"]