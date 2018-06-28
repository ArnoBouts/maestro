FROM debian:jessie
# FROM_DIGEST sha256:d68fe870fe9c7d71d98cf575c42e7c4c3e024c42581eaa116041aa6092430662

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
