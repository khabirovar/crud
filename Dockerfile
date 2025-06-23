FROM golang:alpine as builder

WORKDIR /build

COPY go.mod .

COPY . .

RUN go build -o app main.go

FROM alpine

WORKDIR /

COPY --from=builder /build/app /app

CMD [ "/app" ]