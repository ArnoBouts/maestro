FROM debian:jessie
# FROM_DIGEST sha256:f51cf81db2de8b5e9585300f655549812cdb27c56f8bfb992b8b706378cd517d

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
