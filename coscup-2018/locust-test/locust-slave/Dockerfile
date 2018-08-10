FROM debian:stretch
MAINTAINER fieliapm (fieliapm@gmail.com)

RUN apt-get update \
    && apt-get install -y ca-certificates \
    && apt-get install -y libczmq-dev libzmq3-dev libsodium-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /locust-slave
COPY locust-slave .
EXPOSE 5557

CMD ["sh", "-c", "./locust-slave --master-host=`/sbin/ip route|awk '/default/ { print $3 }'` --master-port=5557 --rpc=zeromq"]
