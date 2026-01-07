# Multi-stage build for service-center
# Supports building multiple binaries: server, client
# Build specific target: docker build --target server -t service-center:server .
# Build all: docker build -t service-center:all .

# FROM golang:1.25.4 AS builder

# # Cache dependencies
# COPY go.mod go.sum ./
# RUN go mod download

# # Copy source
# COPY . .

# Build all binaries
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/server ./cmd/server
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/client ./cmd/client

FROM alpine
WORKDIR /app
COPY ./bin/server /app/server


# # Server image
# FROM gcr.io/distroless/base-debian12 AS server
# COPY --from=builder /bin/server /server
# EXPOSE 50051
# ENTRYPOINT ["/server"]

# All-in-one image (includes both binaries)
# FROM gcr.io/distroless/base-debian12 AS all
# COPY --from=builder /bin/server /server
# COPY --from=builder /bin/client /client
# EXPOSE 50051
# ENTRYPOINT ["/server"]