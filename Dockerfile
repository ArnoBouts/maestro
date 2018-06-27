FROM debian:jessie
# FROM_DIGEST sha256:8ae2506f34500fab08d15ce55b6fd65be34825d7cf8ebc4d6e1f281b234b3446

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
