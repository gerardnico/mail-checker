# syntax=docker/dockerfile:1

# Go toolchain version 1.19.x
FROM golang:1.19 AS build-stage

# Install Go dependencies
# Slash at the end is important when you copy
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download


# Compile to the root of the filesystem
# Slash at the end is important when you copy
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /mail-checker


# Run the tests in the container
FROM build-stage AS run-test-stage
# RUN go test -v ./...
WORKDIR /
COPY --from=build-stage /mail-checker /mail-checker

CMD ["/mail-checker"]
