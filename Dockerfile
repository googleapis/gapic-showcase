FROM golang:1.10

# Setup directory
WORKDIR /go/src/github.com/googleapis/gapic-showcase
COPY . .

# Install dependencies
RUN go get -d -v ./...
RUN go install -v ./...

# Expose port
EXPOSE 8080

# Run the server.
CMD go run /go/src/github.com/googleapis/gapic-showcase/cmd/server/main.go
