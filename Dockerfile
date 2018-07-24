FROM debian:jessie
# FROM_DIGEST sha256:a64a7d8ff7ff87edb78004ef0b159661546d2ddbd82772128b344c90cf8422ab

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
