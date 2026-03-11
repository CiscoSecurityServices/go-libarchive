FROM golang:1.22 AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y libarchive-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy the Go module files
COPY go.mod ./

# Download the Go module dependencies
RUN go mod download

COPY . .

RUN go mod tidy
