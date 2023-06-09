FROM golang:1.19 AS build-stage

WORKDIR /app

COPY . .

RUN GOOS=linux go build -o main-app .

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app

COPY --from=build-stage /app .

# Create a non-root user
RUN adduser --disabled-password --gecos "" appuser
USER appuser

EXPOSE 9090 
CMD ["./main-app"]