FROM scratch
COPY example /usr/bin/arct
ENTRYPOINT ["/usr/bin/arct"]

