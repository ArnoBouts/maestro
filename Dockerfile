FROM debian:jessie
# FROM_DIGEST sha256:fd42763e1dbe2ed3a16f8710b95b015a0db3fd4e7cb713d390aa35832f5acd64

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
