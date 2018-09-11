FROM debian:jessie
# FROM_DIGEST sha256:53ad4744268d3ae3ab11efffe55693e8f1f0d5cb93742fd1d2f26f0a4a84f839

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
