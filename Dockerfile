FROM golang:1.16-alpine

WORKDIR /continuous-release-radar

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /docker-crr

COPY cronjob.sh /continuous-release-radar/

RUN chmod +x cronjob.sh

RUN /continuous-release-radar/cronjob.sh

EXPOSE 8080