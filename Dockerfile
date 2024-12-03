FROM golang:1.23.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# copy rest of the source code
COPY ./ ./

RUN CGO_ENABLED=0 go build -o main cmd/cli/main.go

# use a minimal runtime image
FROM alpine:3.15

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]
