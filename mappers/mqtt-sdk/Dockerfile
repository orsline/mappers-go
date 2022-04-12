FROM alpine:3.14
RUN mkdir -p kubeedge
COPY ./bin/mqtt-sdk kubeedge/bin/
COPY ./res kubeedge/
WORKDIR kubeedge
ENTRYPOINT ["/kubeedge/bin/mqtt-sdk", "--v", "4"]