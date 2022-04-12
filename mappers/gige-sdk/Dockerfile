FROM ubuntu:16.04
WORKDIR /usr/local/bin
COPY ./res /usr/local/res
COPY ./bin /usr/local/bin
COPY ./baumer/Ubuntu-16.04/arm64 /usr/local/bin/cti
ENV GENICAM_GENTL64_PATH /usr/local/bin/cti
RUN apt update && apt install -y libpng16-dev
RUN sed -i '1a /usr/local/bin/genicam/Linux64_ARM' /etc/ld.so.conf && ldconfig
ENTRYPOINT ["./main","--v","4"]

