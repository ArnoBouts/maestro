FROM debian:jessie
# FROM_DIGEST sha256:90f44b88dd8d80bd0fca08c728591fbc43fe36feed3d38428ae1d6d375d96689

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
