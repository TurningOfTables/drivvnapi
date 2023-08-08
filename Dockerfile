# syntax=docker/dockerfile:1

# Build from source
FROM alpine:latest AS build-stage
RUN apk add --no-cache --update go gcc g++

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -o /app/drivvnapi

# Run tests
FROM build-stage AS run-test-stage
RUN go test

# Deploy app binary into lean image
FROM alpine:latest AS build-release-stage

WORKDIR /app

# Copy only the binary, database and readme
COPY --from=build-stage /app/drivvnapi ./
COPY --from=build-stage /app/readme.MD ./
COPY --from=build-stage /app/data/cardata.db ./data/

EXPOSE 8000
CMD ["./drivvnapi", "-r" "-cache"]