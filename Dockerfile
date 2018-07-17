FROM debian:jessie
# FROM_DIGEST sha256:024b0f9f11ab31fd88ff3e31a5cefa5e0dbeb8f3e3d5517a68cee9314072fed1

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
