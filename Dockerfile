# Build the elastic-jupyter-operator binary
FROM golang:1.17.6-alpine as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Download libs first to use docker buildx caching
RUN go mod download
RUN go mod verify

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 go build -a -o elastic-jupyter-operator main.go

# Use distroless as minimal base image to package the elastic-jupyter-operator binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/elastic-jupyter-operator .
USER nonroot:nonroot

ENTRYPOINT ["/elastic-jupyter-operator"]
