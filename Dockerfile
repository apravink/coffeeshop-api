FROM golang:1.8

RUN mkdir -p /app



WORKDIR /app

ADD . /app

Run go get gopkg.in/mgo.v2
Run go get github.com/apravink/coffeeshop-api

RUN go build ./main.go

CMD ["./main"]
