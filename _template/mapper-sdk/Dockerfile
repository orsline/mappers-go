<<<<<<< .mine
FROM alpine:3.14
WORKDIR /
COPY ./res /res
COPY ./bin/main ./bin/main
ENTRYPOINT ["/bin/main"]
CMD ["--v=4"]
=======
FROM alpine:3.14
RUN mkdir -p kubeedge
COPY ./bin/Template kubeedge/bin/
COPY ./res kubeedge/res/
WORKDIR kubeedge/bin
ENTRYPOINT ["./Template"]
>>>>>>> .theirs
