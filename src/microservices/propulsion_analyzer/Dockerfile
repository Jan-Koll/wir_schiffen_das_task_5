# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.22.4 AS build-stage

WORKDIR /app

COPY *.mod *.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /propulsion_analyzer

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /propulsion_analyzer /propulsion_analyzer

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/propulsion_analyzer"]