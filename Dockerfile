FROM golang:1.9 AS build
ADD . /go/src/github.com/albertogviana/docker-swarm-service-status
WORKDIR /go/src/github.com/albertogviana/docker-swarm-service-status
RUN go get -d -v -t
RUN CGO_ENABLED=0 GOOS=linux go build -v -o docker-swarm-service-status


FROM alpine
LABEL maintaner "Alberto Guimaraes Viana <viana@vesseltracker.com>"

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY --from=build /go/src/github.com/albertogviana/docker-swarm-service-status/docker-swarm-service-status /usr/local/bin/docker-swarm-service-status

EXPOSE 8080
CMD ["docker-swarm-service-status"]

HEALTHCHECK --interval=5s --start-period=3s --timeout=5s CMD wget -qO- "http://localhost:8080/v1/docker-swarm-service-status/health"

RUN chmod +x /usr/local/bin/docker-swarm-service-status