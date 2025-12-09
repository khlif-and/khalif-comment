FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

FROM alpine:latest

# Menghapus ffmpeg karena service ini hanya teks
RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Jakarta

WORKDIR /root/

COPY --from=builder /app/server .

# Port diubah ke 8083 agar tidak bentrok dengan stories (8082)
EXPOSE 8083

CMD ["./server"]