FROM debian:jessie
# FROM_DIGEST sha256:58c95ab3ce7069ee18cac3cc2fcb8a1ab2cdcbcc2e7ad2aed851bb888f9e8f6d

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
