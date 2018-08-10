FROM debian:latest
#FROM golang:latest
MAINTAINER fieliapm (fieliapm@gmail.com)

RUN apt-get update \
    && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/*

#WORKDIR /go/src/locust-test/locust-slave
#COPY locust-slave.go .
#RUN go get -v
#RUN go build -v -o locust-slave locust-slave.go
#RUN mkdir /locust-slave && cp locust-slave /locust-slave

WORKDIR /locust-slave
COPY locust-slave .
EXPOSE 5557

# for Linux
CMD ["sh", "-c", "./locust-slave --master-host=`/sbin/ip route|awk '/default/ { print $3 }'` --master-port=5557 --rpc=zeromq"]
# for MacOSX
#CMD ["sh", "-c", "./locust-slave --master-host=host.docker.internal --master-port=5557 --rpc=zeromq"]
