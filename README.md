# Golang DB Users Platform

This is a custom platform for multiple databases users administration.

# Requirements

    - Go
    - Docker

# Getting Started

    go mod download
    go run .


# Docker Image 
## Enviroment variables
    - CLERK=clerk api key
    - DB_NAME=db name
    - DB_USER=db user
    - DB_PASSWORD=db password
    - DB_ADDRESS=db address

sample:
```
    User:   "root",
    Passwd: "test",
    Net:    "tcp",
    Addr:   "127.0.0.1:3001",
    DBName: "db_platform",
```

To run the test databases go to the sample folder and run: `docker compose up`

# Useful commands

To shel into any of the databases use: `docker exec -u "user" -it "container name" "shell(Example: mysql, psql)" -pPasssword`

To build an image of this application use: `docker build --tag "image name" .`

To run the image for this application use: `docker run --rm --name "container name" -e "clerk=clerk key" -p 3000:8080/tcp  "image name" `

Tailwind use: `npx tailwindcss -i ./src/server/assets/css/styles.css -o ./src/server/assets/css/output.css --watch`
