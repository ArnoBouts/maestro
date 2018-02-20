FROM debian:jessie
# FROM_DIGEST sha256:23916ac1bc20307fe4c54938c559efc69f35e4431c5d44003dc9d6c777c91e48

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
