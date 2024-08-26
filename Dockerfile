FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -tags prod -o db-platform


FROM registry.access.redhat.com/ubi9-minimal:latest

WORKDIR /app

# COPY --from=builder /app/src/web ./src/web
COPY --from=builder /app/db-platform .

EXPOSE 3000/tcp

CMD ["./db-platform"]