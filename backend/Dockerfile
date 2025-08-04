FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o data-ingestion-backend ./cmd

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/data-ingestion-backend .

EXPOSE ${PORT}

CMD ["sh", "-c", "./data-ingestion-backend"]