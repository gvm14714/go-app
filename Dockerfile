FROM golang:1.19 AS build-stage

WORKDIR /app

COPY ./app .

RUN GOOS=linux go build -o main-app .

# Deploy the application binary into a lean image
FROM debian:11-slim AS build-release-stage

RUN apt-get update
RUN apt-get install -y telnet
WORKDIR /app

COPY --from=build-stage /app .

# Create a non-root user
#RUN adduser --disabled-password --gecos "" appuser
#USER appuser

EXPOSE 9090 
CMD ["./main-app"]