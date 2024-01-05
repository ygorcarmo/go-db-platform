# Golang DB Users Platform

This is a custom platform for multiple databases users administration.

# Requirements

    - Go
    - Docker

# Getting Started

    go mod download
    go run .

To run the test databases go to the sample folder and run: `docker compose up`

# Useful commands

To shel into any of the databases use: `docker exec -u "user" -it "container name" "shell(Example: mysql, psql)" -pPasssword`

To build an image of this application use: `docker build --tag "image name" .`

To run the image for this application use: `docker run --rm --name "container name" -e "clerk=clerk key" -p 3000:8080/tcp  "image name" `
