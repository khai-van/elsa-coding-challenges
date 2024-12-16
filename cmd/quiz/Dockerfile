FROM golang:1.22 as builder

WORKDIR /app

RUN --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download

RUN --mount=type=bind,target=. \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/service ./cmd/quiz

FROM gcr.io/distroless/static as runner

WORKDIR /go/bin
COPY ./configs/ ./configs/
COPY --from=builder /bin/service service

# Execute app
ENTRYPOINT ["./service"]