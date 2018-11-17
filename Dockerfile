FROM debian:jessie
# FROM_DIGEST sha256:14e15b63bf3c26dac4f6e782dbb4c9877fb88d7d5978d202cb64065b1e01a88b

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
