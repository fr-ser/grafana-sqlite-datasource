FROM golang:1.14

RUN apt-get update -y
RUN apt-get install -y \
    gcc-aarch64-linux-gnu \
    libc6-dev-arm64-cross \
    gcc-mingw-w64-x86-64

