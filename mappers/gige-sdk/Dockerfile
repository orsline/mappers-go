FROM alpine:3.14
WORKDIR /
COPY ./res ./res
COPY ./bin ./bin
COPY ./baumer/Ubuntu-16.04/arm64 ./bin/cti
ENV GENICAM_GENTL64_PATH /bin/cti
RUN sed -i '1a /bin/genicam/Linux64_ARM' /etc/ld.so.conf && ldconfig
RUN sysctl -w net.core.rmem_max=33554432 && sysctl -w net.core.netdev_max_backlog=2000 && sysctl -w net.core.netdev_budget=600
ENTRYPOINT ["/bin/main"]
CMD ["--v=4"]
EXPOSE 80/tcp
EXPOSE 80/udp

