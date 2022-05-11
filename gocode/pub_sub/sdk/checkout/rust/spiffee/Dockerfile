FROM ubuntu:20.04

ARG DEBIAN_FRONTEND=noninteractive
RUN apt update -y && apt-get upgrade -y
RUN apt-get install -y curl gcc g++ git make pkg-config cmake libssl-dev protobuf-compiler

ADD ./spiffee .

CMD ./spiffee