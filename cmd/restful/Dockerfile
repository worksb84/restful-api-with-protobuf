FROM golang:1.23.2-alpine

RUN mkdir /AAAA
WORKDIR /AAAA
COPY . .
RUN go mod download

RUN go build -o AAAA ./

ENV GIN_MODE release

EXPOSE 8001
CMD [ "./AAAA" ]
