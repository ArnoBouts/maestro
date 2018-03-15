FROM debian:jessie
# FROM_DIGEST sha256:44d53bda18cc6b564cedca4bd4d7cfacce37a3f43a439ea8a73c4d8c32f53400

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
