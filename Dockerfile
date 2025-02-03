FROM golang:1.21-alpine

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY main.go ./

COPY input/ /app/input/


RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.19

WORKDIR /app

COPY --from=0 /app/main .
COPY --from=0 /app/input /app/input

EXPOSE 8080

CMD ["./main"] 