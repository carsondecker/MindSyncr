FROM golang:1.25

WORKDIR /usr/src/app

# Copy workspace files first
COPY go.work ./

# Copy all modules (so workspace works)
COPY utils ./utils
COPY api-gateway ./api-gateway
COPY ws-service ./ws-service

WORKDIR /usr/src/app/ws-service

# Cache dependencies
COPY ws-service/go.mod ws-service/go.sum ./
RUN go mod download

# Build
RUN go build -v -o /usr/local/bin/app ./cmd/MindSyncr-WS

CMD ["app"]