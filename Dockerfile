FROM golang:1.20.4-bullseye as builder

WORKDIR /app

COPY . .

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -a -o main .

FROM scratch

COPY --from=builder /app/main .

CMD ["./main"]
