FROM golang:1.17.9-alpine3.14
RUN mkdir /app

ADD . /app
WORKDIR /app
RUN GO111MODULE=on go mod download
RUN go build -o main .

CMD ["/app/main"]