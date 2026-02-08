FROM golang:1.23

WORKDIR /app
COPY . . 
#go get -d -v ./... - скачивает и обновляет зависимости, go mod dowland - скачивает зависимости
RUN cd server && go get -d -v ./...
RUN cd server && go build -o ../application ./src/app/main.go
EXPOSE 80
CMD ["./application"]
