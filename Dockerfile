FROM golang:1.23-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/go-link-checker ./main.go

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/go-link-checker /bin/go-link-checker
RUN chmod +x /bin/go-link-checker

ENTRYPOINT ["/bin/go-link-checker"]