# Use the official Golang base image
FROM golang:1.19-alpine as builder

# Set the working directory
WORKDIR /app

RUN apk add --no-cache sqlite

RUN apk add build-base

RUN apk --no-cache add make git gcc libtool musl-dev ca-certificates dumb-init 


# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

COPY . .

RUN go mod download

# # Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o xgo-service ./main/src/cmd/main.go 

# # Use a minimal base image
FROM alpine:latest

# # Set the working directory
WORKDIR /app

# # Copy the built binary from the builder stage
COPY --from=builder /app/xgo-service .

COPY xgo.ini /app

RUN apk update && apk add bash

# RUN apk add librdkafka-dev pkgconf
RUN apk update && apk upgrade

RUN apk add --no-cache sqlite

# # Expose the desired port
EXPOSE 8100

# # Run the application
CMD ["./xgo-service"]
