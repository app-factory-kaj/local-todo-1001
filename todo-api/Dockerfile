FROM golang:1.25-alpine AS builder
WORKDIR /src
# Stdlib-only service: no go.sum to copy or download against.
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -ldflags='-s -w' -o /out/app ./

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/app /app
EXPOSE 9090
ENTRYPOINT ["/app"]
