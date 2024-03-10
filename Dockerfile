FROM scratch
COPY gwm /bin/gwm
ENTRYPOINT ["/bin/gwm"]
