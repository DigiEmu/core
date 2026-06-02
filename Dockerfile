## Multi-stage build: compile the digiemu CLI and produce a minimal runtime image
FROM golang:1.25.7-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /src/cmd/digiemu
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o /usr/local/bin/digiemu .

FROM scratch AS runtime
COPY --from=builder /usr/local/bin/digiemu /usr/local/bin/digiemu
COPY --from=builder /src/testdata/core_2_conformance /opt/testdata/core_2_conformance
ENTRYPOINT ["/usr/local/bin/digiemu"]
CMD ["experimental", "conformance", "run", "/opt/testdata/core_2_conformance"]
