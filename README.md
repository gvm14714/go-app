# go-app-intern

go-app-intern is a web app written in Go that is dockerized, automated using jenkins and deployed on minikube.

Contains:
- Dockerfile
- dockercompose
- Jenkinsfile
- kubernetes files

## Prerequistis

- Docker
- docker compose
- Jenkins
- helm
- Minikube

## How to build this app

```bash
go build -o goviolin .
```
This will produce an artifiact called goviolin then to run it
```bash
./main-app
```
![image](https://github.com/ahmedelmelegy/GoViolin/assets/62904201/5db629ec-02ca-4d61-bc47-557200030ab5)
## Dockerize application
### Single-Stage Dockerfile

```bash
FROM golang:1.19

WORKDIR /main-app

COPY ./app . RUN go mod init RUN go build -o app .

EXPOSE 9090
CMD ["./main-app"]
```
To build image
```bash
docker build . -t app
```
![image](https://github.com/gAhmed-Elmelegy/go-app-intern/assets/136341359/d9ee8997-8d39-43fc-9e3a-d0b02d0ded6e)
The size of the image is more than 1 GB!!!
### Multi-Stage Dockerfile
As we saw single stage Dockerfile generated an image that has very big size and that isn`t effiecient
```bash
FROM golang:1.19 AS build-stage

WORKDIR /app

COPY ./app .

RUN GOOS=linux go build -o main-app .

# Deploy the application binary into a lean image
FROM debian:11-slim AS build-release-stage

WORKDIR /app

COPY --from=build-stage /app .

#Create a non-root user
RUN adduser --disabled-password --gecos "" appuser
USER appuser

EXPOSE 9090 
CMD ["./main-app"]
```
To build image
```bash
docker build . -t ahmedelmelegy3570/app-multistage
```
What we benifted from that?
The size of the image is greatly reduced
![image](https://github.com/gAhmed-Elmelegy/go-app-intern/assets/136341359/1d839a6b-504b-4ac6-92ed-021e78978b1e)
Now it is less than 0.1 GB!

To run container with the built image, we will map port 9090 from host to port 9090 in the container
```bash
docker run -p 9090:9090 ahmedelmelegy3570/app-multistage
```
### To make the image more secure
to create a user and use it instead of root user
```bash
RUN adduser --disabled-password --gecos "" appuser
USER appuser
```
## Docker-compose
in dockercompose.yaml, I combined both app and mysql db so they have the same network and I made the app depend on mysql
```bash
version: '3'
services:
  app:
    image: app-multistage
    ports:
      - 9090:9090
    depends_on:
      - mysql
  mysql:
    image: mysql:5.7
    environment:
      - MYSQL_ROOT_PASSWORD=1234
      - MYSQL_DATABASE=mydb
      - MYSQL_USER=ahmed
      - MYSQL_PASSWORD=1234
    volumes:
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - 3306:3306
    command: ["--log-error-verbosity=3"]
```
To apply this docker compose file
```bash
docker-compose up
```
![image](https://github.com/gAhmed-Elmelegy/go-app-intern/assets/136341359/cf3d91b0-adf3-4407-a804-459fe8517aa1)

