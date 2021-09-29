# Build the binary
FROM golang:1.17 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY server.go server.go

COPY db/ db/
COPY lib/ lib/
COPY proto/ proto/
COPY static/ static/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a server.go

# Use distroless as minimal base image to package the binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /

COPY --from=builder /workspace/server .
COPY --from=builder /workspace/static/ ./static/
COPY --from=builder /workspace/db/ ./db/
USER 1324:1324

ENTRYPOINT ["/server"]
