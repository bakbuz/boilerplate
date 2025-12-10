# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Static binary build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server cmd/server/main.go

# Stage 2: Minimal Runtime (Opsiyonel, eğer container ile deploy edeceksen)
FROM scratch
COPY --from=builder /app/server /server
# SSL sertifikaları Supabase bağlantısı için gereklidir
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 50051
ENTRYPOINT ["/server"]