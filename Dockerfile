FROM golang:1.9 AS build
ADD . /go/src/github.com/albertogviana/docker-swarm-deployment-status
WORKDIR /go/src/github.com/albertogviana/docker-swarm-deployment-status
RUN go get -d -v -t
RUN CGO_ENABLED=0 GOOS=linux go build -v -o docker-swarm-deployment-status


FROM alpine
LABEL maintaner "Alberto Guimaraes Viana <viana@vesseltracker.com>"

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY --from=build /go/src/github.com/albertogviana/docker-swarm-deployment-status/docker-swarm-deployment-status /usr/local/bin/docker-swarm-deployment-status

EXPOSE 8080
CMD ["docker-swarm-deployment-status"]

RUN chmod +x /usr/local/bin/docker-swarm-deployment-status