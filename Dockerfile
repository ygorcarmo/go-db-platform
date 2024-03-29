FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o db-platform


FROM alpine:latest

WORKDIR /app

# COPY --from=builder /app/src/web ./src/web
COPY --from=builder /app/db-platform .

EXPOSE 3000/tcp

CMD ["./db-platform"]