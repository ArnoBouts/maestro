FROM debian:jessie
# FROM_DIGEST sha256:e4662a38a8504fae8324e089cd95e1b33cc45ab876ca699969e1e92ced74bd13

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
