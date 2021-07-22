FROM scratch
COPY arct /usr/bin/arct
ENTRYPOINT ["/usr/bin/arct"]

