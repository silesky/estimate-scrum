FROM golang:latest
RUN mkdir /app
ADD .
WORKDIR /app
RUN go get -d ./...
RUN go build -o main .
CMD ["/app/main"]
