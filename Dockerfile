FROM golang:1.11

# Setup directory
WORKDIR /go/src/github.com/googleapis/gapic-showcase
COPY . .

# Install showcase
RUN ["go", "get"]
RUN ["go", "install"]

# Expose port
EXPOSE 7469

# Run the server.
ENTRYPOINT ["gapic-showcase"]
CMD ["start"]
