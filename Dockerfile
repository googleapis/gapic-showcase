FROM golang:1.11-alpine AS builder

# Install git and gcc.
RUN apk add --no-cache git gcc musl-dev

# Setup directory.
WORKDIR /go/src/github.com/googleapis/gapic-showcase
COPY . .

# Compile for Linux.
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64

# Install showcase.
RUN go get ./...
RUN go build -installsuffix cgo \
  -ldflags="-w -s" \
  -o /go/bin/gapic-showcase \
  ./cmd/gapic-showcase

# Start a fresh image, and only copy the built binary.
FROM scratch
COPY --from=builder /go/bin/gapic-showcase /go/bin/gapic-showcase

# Expose port
EXPOSE 7469

# Run the server.
ENTRYPOINT ["/go/bin/gapic-showcase"]
CMD ["run"]
