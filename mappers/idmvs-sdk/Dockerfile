FROM alpine:3.14
WORKDIR /
COPY ./res /res
COPY ./bin/main ./bin/main
ENTRYPOINT ["/bin/main"]
CMD ["--v=4"]