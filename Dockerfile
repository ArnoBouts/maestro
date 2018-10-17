FROM debian:jessie
# FROM_DIGEST sha256:d671cef731e199f6e1796798ce1067c690d05036a7040ad1349bd3c44024dc8b

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
