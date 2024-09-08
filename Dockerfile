FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/http/main.go

FROM alpine:latest

RUN apk --no-cache add libc6-compat
RUN apk add vips

WORKDIR /root/

COPY --from=builder /app/main .

RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]