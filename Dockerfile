FROM golang:1.21

COPY . /app

WORKDIR /app/cmd

RUN go build main.go

CMD ["./main"]