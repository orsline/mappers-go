FROM alpine:3.14
RUN mkdir -p kubeedge
COPY ./bin/idmvs-sdk kubeedge/bin/
COPY ./res kubeedge/
WORKDIR kubeedge
ENTRYPOINT ["/kubeedge/bin/idmvs-sdk", "--v", "4"]