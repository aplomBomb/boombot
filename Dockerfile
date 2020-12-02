FROM golang:alpine

RUN apk update && apk add --no-cache git && apk add ffmpeg

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN mkdir /boombot

ADD . /boombot

WORKDIR /boombot

RUN go build -o boombot 

RUN chmod +x start.sh

CMD [ "/boombot/boombot" ] 
