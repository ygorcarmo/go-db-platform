FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-db-platform

EXPOSE 8080/tcp

CMD [ "/docker-db-platform" ]