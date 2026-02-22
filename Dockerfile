FROM golang:alpine AS builder

ARG VERSION="dev"
ARG GitCommit="Unknown"
ARG BuildTime="Unknown"

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    go get github.com/gofiber/fiber/v2@v2.52.0 \
    go get github.com/gofiber/fiber/v2/middleware/basicauth@v2.52.0 \
    go get github.com/gofiber/fiber/v2/middleware/cors@v2.52.0 \
    go get github.com/gofiber/fiber/v2/middleware/logger@v2.52.0 \
    go get github.com/gofiber/fiber/v2/middleware/recover@v2.52.0
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 \
    go build -trimpath \
    -ldflags=" \
    -s -w \
    -X 'github.com/kiss2u/SaveAny-Bot/config.Version=${VERSION}' \
    -X 'github.com/kiss2u/SaveAny-Bot/config.GitCommit=${GitCommit}' \
    -X 'github.com/kiss2u/SaveAny-Bot/config.BuildTime=${BuildTime}' \
    -X 'github.com/kiss2u/SaveAny-Bot/config.Docker=true' \
    " \
    -o saveany-bot .

FROM alpine:latest

RUN apk add --no-cache curl ffmpeg yt-dlp

WORKDIR /app

COPY --from=builder /app/saveany-bot .
COPY entrypoint.sh .

RUN chmod +x /app/saveany-bot && \
    chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
