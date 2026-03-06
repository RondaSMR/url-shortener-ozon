FROM golang:1.25-alpine AS builder

WORKDIR /src

ARG TARGETOS
ARG TARGETARCH

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -trimpath -ldflags="-s -w" -o /out/url-shortener ./cmd/app

FROM alpine:3.20

RUN adduser -D -H -s /sbin/nologin appuser
USER appuser

COPY --from=builder /out/url-shortener /usr/local/bin/url-shortener

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/url-shortener"]
