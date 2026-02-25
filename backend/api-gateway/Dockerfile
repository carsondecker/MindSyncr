FROM golang:1.25

WORKDIR /usr/src/app

# Copy workspace files first
COPY go.work ./

# Copy all modules (so workspace works)
COPY utils ./utils
COPY ws-service ./ws-service
COPY api-gateway ./api-gateway

WORKDIR /usr/src/app/api-gateway

# Cache dependencies
COPY api-gateway/go.mod api-gateway/go.sum ./
RUN go mod download

# Build
RUN go build -v -o /usr/local/bin/app ./cmd/MindSyncr

CMD ["app"]