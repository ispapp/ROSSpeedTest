FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./sptest ./main.go
 
 
FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/runner .
EXPOSE 8080
ENTRYPOINT ["./sptest"]