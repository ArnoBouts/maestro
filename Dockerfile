FROM debian:jessie
# FROM_DIGEST sha256:c345dbca9a4184c1d717e6fd26674218e1437775dba1261039c1fcce9fbb6d1b

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
