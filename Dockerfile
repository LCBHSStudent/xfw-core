FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive

RUN mkdir -p /xfw/xfw-core && cd /xfw &&\
    apt-get update && apt-get install -y mysql-server wget tar git &&\
    wget https://github.com/LCBHSStudent/xfw-core/releases/download/v1.0.3/xfw-core-v1.0.3.tar.gz &&\
    tar -zxvf ./xfw-core-v1.0.3.tar.gz -C ./xfw-core && cd /xfw &&\
    git clone https://github.com/sheepzh/poetry && mkdir -p ./xfw-core/share/poem && mv ./poetry/data/* ./xfw-core/share/poem &&\
    service mysql start &&\
    echo "create database xfw;\
    create database homospace;\
    CREATE USER 'dfw'@'%' IDENTIFIED BY '^52]pt*xz+g^03_C#YHb';\
    grant all privileges on *.* to 'dfw'@'%';\
    flush privileges;" | mysql -h 127.0.0.1 -u$(awk 'NR==4{print $3}' /etc/mysql/debian.cnf) -p$(awk 'NR==5{print $3}' /etc/mysql/debian.cnf) &&\
    service mysql stop && apt-get remove -y git && apt autoremove -y && apt-get clean

WORKDIR /xfw

CMD "xfw-core/scripts/run.sh"
