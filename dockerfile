
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/forum ./cmd

FROM alpine:3.18
RUN addgroup -S app && adduser -S -G app app
COPY --from=builder /app/forum /usr/local/bin/forum
COPY --from=builder /src/templates /app/templates
COPY --from=builder /src/static /app/static
COPY --from=builder /src/upload /app/upload
COPY --from=builder /src/database/schema.sql /app/schema.sql

RUN chown -R app:app /app /usr/local/bin/forum
USER app
WORKDIR /app

ENV PORT=8080


EXPOSE 8080

CMD ["/usr/local/bin/forum"]