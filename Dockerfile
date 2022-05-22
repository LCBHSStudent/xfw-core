FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update &&\
    apt-get install -y git mysql-server software-properties-common --fix-missing &&\
    add-apt-repository -y ppa:longsleep/golang-backports && apt-get update && apt-get install -y golang &&\
    mkdir -p /xfw && cd /xfw &&\
    git clone https://github.com/LCBHSStudent/xfw-core &&\
    mkdir -p xfw-core/bin && cd xfw-core && go build -o bin/xfw-core ./src/core && cd bin && chmod +x go-cqhttp && cd /xfw &&\
    git clone https://github.com/sheepzh/poetry && mkdir -p ./xfw-core/share/poem && mv ./poetry/data/* ./xfw-core/share/poem &&\
    service mysql start &&\
    echo "create database xfw;\
    CREATE USER 'dfw'@'%' IDENTIFIED BY '^52]pt*xz+g^03_C#YHb';\
    grant all privileges on *.* to 'dfw'@'%';\
    flush privileges;" | mysql -h 127.0.0.1 -u$(awk 'NR==4{print $3}' /etc/mysql/debian.cnf) -p$(awk 'NR==5{print $3}' /etc/mysql/debian.cnf) &&\
    service mysql stop && apt-get remove -y golang software-properties-common git && apt autoremove -y && apt-get clean

WORKDIR /xfw

CMD service mysql start && /bin/bash
