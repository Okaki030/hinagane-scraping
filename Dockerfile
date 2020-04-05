FROM ubuntu:latest

RUN apt-get update \
    && apt-get install -y wget git gcc-7 g++-7 build-essential file \
    && wget -P /tmp https://dl.google.com/go/go1.12.6.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf /tmp/go1.12.6.linux-amd64.tar.gz \
    && rm /tmp/go1.12.6.linux-amd64.tar.gz

ENV GOPATH $HOME/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR $GOPATH/src

ENV GO111MODULE=on
ENV CGO_LDFLAGS="-L/path/to/lib -lmecab -lstdc++"
ENV CGO_CFLAGS="-I/path/to/include"

RUN apt-get update \
    && apt-get install -y curl sudo cron swig \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /go
RUN git clone https://github.com/taku910/mecab.git
WORKDIR /go/mecab/mecab
RUN ./configure  --enable-utf8-only \
    && ls \
    && make \
    && make check \
    && make install \
    && ldconfig

WORKDIR /go/mecab/mecab-ipadic
RUN ./configure --with-charset=utf8 \
    && make \
    &&make install

WORKDIR /go
RUN git clone --depth 1 https://github.com/neologd/mecab-ipadic-neologd.git
WORKDIR /go/mecab-ipadic-neologd
RUN ./bin/install-mecab-ipadic-neologd -n -y

WORKDIR /go/src
EXPOSE 8080
CMD ["go", "run", "/go/src/cmd/main.go"]