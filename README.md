# DB Users Platform

This is a custom platform for multiple databases users administration.

# Requirements to run it locally

    - Go
    - Docker

# Getting Started

    go mod download
    make dev

# Docker Image

## Enviroment variables

    - LISTEN_ADDR=":3002"
    - DB_USER="apt_db_platform"
    - DB_PASSWORD="CHANGEME"
    - DB_ADDRESS="localhost:3001"
    - DB_NAME="db_platform"
    - JWT_SECRET="DONOTSHOWTHISTOANYONE"

To run the test databases go to the sample folder and run: `docker compose up`

# Useful commands

To shell into any of the databases use: `docker exec -u "user" -it "container name" "shell(Example: mysql, psql)" -pPasssword`

To build an image of this application use: `docker build --tag "image name" .`

To run the image for this application use: `docker run --rm --name "container name" -e "ADD ALL ENVIROMENT VARS HERE" -p 3000:8080/tcp  "image name" `

Tailwind use: `make css`
