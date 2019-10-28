FROM golang:latest
MAINTAINER Marshall Lee <marshall.s.lee@gmail.com>

RUN mkdir /app

ADD cmd/exec /app/cmd/exec

WORKDIR /app/cmd/exec
COPY . .

RUN go get -t github.com/gorilla/mux
RUN go get -t github.com/googollee/go-socket.io
RUN go get -t github.com/marshallslee/chitchat-push/handler

EXPOSE 12379

RUN go build -o socketiosample .

CMD ./socketiosample