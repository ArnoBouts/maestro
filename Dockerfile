FROM debian:jessie
# FROM_DIGEST sha256:a69efd2cf1f83d7fe6a374aaf27974ed031d72b2886acb4b730c6a2017db0c3d

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
