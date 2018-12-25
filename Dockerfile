FROM golang:1.11 as builder
WORKDIR /go/src/github.com/alexmorten/mhist
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o mhist main/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir app
COPY --from=builder /go/src/github.com/alexmorten/mhist/mhist /app
WORKDIR /app
CMD ["./mhist"]
EXPOSE 6666 6667
