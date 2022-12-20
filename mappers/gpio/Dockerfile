FROM ubuntu:20.04
WORKDIR /usr/local/bin
COPY ./res /res
COPY ./bin/main ./usr/local/bin/gpio
ENTRYPOINT ["./gpio","--v","4"]
