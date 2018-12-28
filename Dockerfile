FROM golang:1.11-alpine

# Install git
RUN apk add --no-cache git gcc musl-dev

# Setup directory
WORKDIR /go/src/github.com/googleapis/gapic-showcase
COPY . .

# Use go modules.
ENV GO111MODULE on

# Install showcase
RUN ["go", "get"]
RUN ["go", "install"]

# Expose port
EXPOSE 7469

# Run the server.
ENTRYPOINT ["gapic-showcase"]
CMD ["run"]
