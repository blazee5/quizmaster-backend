FROM golang:1.21

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

WORKDIR /app/cmd

RUN go build main.go

CMD ["./main"]