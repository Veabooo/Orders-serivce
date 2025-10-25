FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o order-service .
RUN go build -o publisher ./Publish/publish.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/order-service .
COPY --from=builder /app/publisher .
COPY templates ./templates
COPY model.json .
EXPOSE 8080
CMD ["./order-service"]