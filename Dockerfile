FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0

RUN go build -o main ./cmd/app

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]