FROM golang:1.25.1
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o app main.go
CMD ["./app"]
