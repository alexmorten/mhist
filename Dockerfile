FROM golang:1.13 as builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o mhist main/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir app
COPY --from=builder /build/mhist /app
WORKDIR /app
CMD ["./mhist"]
EXPOSE 6666
