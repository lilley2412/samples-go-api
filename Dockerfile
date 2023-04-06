FROM golang:1.19 as builder

WORKDIR /app

ENV HOME=/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/main

FROM alpine

COPY --from=builder /app/main /app/main

ENTRYPOINT [ "/app/main" ]