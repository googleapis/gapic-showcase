FROM golang:1.11-alpine AS builder

# Install git and gcc.
RUN apk add --no-cache git gcc musl-dev

# Setup directory.
WORKDIR /go/src/github.com/googleapis/gapic-showcase
COPY . .

# Use go modules, and only compile for Linux.
ENV GO111MODULE on
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64

# Install showcase.
RUN go get
RUN go build -a \
  -installsuffix cgo \
  -ldflags="-w -s" \
  -o /go/bin/gapic-showcase

# Start a fresh image, and only copy the built binary.
FROM scratch
COPY --from=builder /go/bin/gapic-showcase /go/bin/gapic-showcase

# Expose port
EXPOSE 7469

# Run the server.
ENTRYPOINT ["/go/bin/gapic-showcase"]
CMD ["run"]
