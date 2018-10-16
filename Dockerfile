FROM debian:jessie
# FROM_DIGEST sha256:3a50c98172c3a0571334bdca09b263e3b525002f906f9fcd3b8fa831d85418a2

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
