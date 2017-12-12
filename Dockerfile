FROM debian:jessie
# FROM_DIGEST sha256:1ed4d5996abb15f31db8e31e300f3327722d2898664d3f89bd522b3636e6b165

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
