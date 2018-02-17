FROM debian:jessie
# FROM_DIGEST sha256:aa7a65b8796ea9e260c485f810b39929890d6ea630484dc75aa51a32c93af8ce

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
