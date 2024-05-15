FROM golang:1.21 AS builder

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -v -o bin/api feather-proxy/cmd

FROM debian:trixie AS runner

WORKDIR /go/bin

EXPOSE 8000

COPY --from=builder /go/src/app/bin/api /go/bin/api

CMD ["./api"]

