FROM alpine:3.14
RUN mkdir -p kubeedge
COPY ./bin/gpio-sdk kubeedge/bin/
COPY ./res kubeedge/res/
WORKDIR kubeedge/bin
ENTRYPOINT ["./gpio-sdk", "--v", "4"]
