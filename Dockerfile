FROM debian:jessie
# FROM_DIGEST sha256:d74dea994f22f51a5b39ecd1501ed285f307c285f8975d9373d8cc1d28e9e847

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
