FROM debian:jessie
# FROM_DIGEST sha256:f29d0c98d94d6b2169c740d498091a9a8545fabfa37f2072b43a4361c10064fc

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
