FROM golang:latest

COPY . .

WORKDIR .

RUN go build main.go

ENTRYPOINT ./main
