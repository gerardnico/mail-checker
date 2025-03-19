# syntax=docker/dockerfile:1
# https://docs.docker.com/guides/golang/build-images/

# Go toolchain version 1.23.x
FROM golang:1.23 AS build-stage

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
#FROM build-stage AS run-test-stage
#RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /
COPY --from=build-stage /mail-checker /mail-checker

CMD ["/mail-checker"]
