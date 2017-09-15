FROM debian:jessie
# FROM_DIGEST sha256:f31c837c5dda4c4730a6f1189c9ae39031e966410d7b0885863bd6c43b5ff281

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
