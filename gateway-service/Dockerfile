#build stage
FROM golang AS builder
RUN apt update && apt install git
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main src/cmd/main/main.go

#final stage
FROM debian:latest
RUN apt update && apt install -y ca-certificates
COPY --from=builder app app
WORKDIR /app
