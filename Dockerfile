FROM golang:1.19 AS build-stage

WORKDIR /app

COPY . .

RUN GOOS=linux go build -o main .

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app

COPY --from=build-stage /app .

EXPOSE 9090 
CMD ["./main"]