FROM golang:1.23

WORKDIR /app
COPY . . 
RUN go get -d -v ./...
RUN go build -o application cmd/app/main.go
EXPOSE 80
CMD ["./application"]
