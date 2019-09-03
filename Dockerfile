##### compile backend in official golang image #####
FROM golang:1.12 as builder
WORKDIR /app
COPY backend ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o feed-squirrel .

##### package binary into alpine image for distribution ####
FROM alpine:3.10
WORKDIR /root/
COPY --from=builder /app/feed-squirrel .
COPY frontend/build frontend
EXPOSE 8080
CMD ["./feed-squirrel"]
