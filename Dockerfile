FROM scratch

COPY maestro /

EXPOSE 80

ENTRYPOINT ["/maestro"]
