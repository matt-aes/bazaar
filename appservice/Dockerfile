FROM golang:latest as build-env

# All these steps will be cached
RUN mkdir /app
WORKDIR /app

# Modules: copy go.mod, go.sum
COPY go.mod .

# NOTE: go.sum not necessary if no modules exist
# COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the source directory (including all subdirectories) to the work directory
COPY . .

# Build the binary
RUN go build -o /go/bin/application .

# Expose our service port
EXPOSE 8080
ENTRYPOINT ["/go/bin/application"]