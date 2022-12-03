# syntax=docker/dockerfile:1

##
## Stage 1: Build the application
##
FROM golang:1.19.3-alpine as builder

LABEL maintainer="Dipak Parmar <hi@dipak.tech>"

WORKDIR /app

# Copy go mod and dependencies to build
COPY go.mod ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -o server .

##
## Stage 2: Deploy the application
##
FROM alpine:3.17.0

WORKDIR /app

COPY --from=builder /app/server .

ENV DOMAIN=example.com

# Expose port 8080 to the outside world
EXPOSE 8080

ENTRYPOINT [ "/app/server", "start" ]

