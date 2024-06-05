FROM golang:1.22.2

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main ./cmd/server/main.go

RUN chmod +x main

EXPOSE 8080

CMD [ "./main" ]

