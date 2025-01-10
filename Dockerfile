FROM golang:1.23.4-alpine AS builder

WORKDIR /build
COPY . .

RUN go mod download
RUN go build -o ./url-shortener

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /build/url-shortener ./url-shortener
CMD ["/app/url-shortener"]