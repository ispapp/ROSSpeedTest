FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -tags=jsoniter -o runner .
 
 
FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/runner .
EXPOSE 8080
ENTRYPOINT ["./runner"]