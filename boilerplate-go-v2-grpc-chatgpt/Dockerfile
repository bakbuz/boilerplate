FROM golang:1.25.5 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /bin/shop ./cmd/server

FROM scratch
COPY --from=builder /bin/shop /bin/shop
EXPOSE 50051
ENTRYPOINT ["/bin/shop"]
