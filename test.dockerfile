FROM golang:latest AS builder

WORKDIR /workdir
COPY go.* .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o test-worker ./test

FROM scratch
COPY --from=builder /workdir/test-worker /
CMD ["/test-worker"]
