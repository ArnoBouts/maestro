FROM debian:jessie
# FROM_DIGEST sha256:d4d4bc28ac378d0890cc784873d1d8eeec2f8c331e15d204ddd7aca8d18bce84

WORKDIR /maestro

COPY maestro /maestro/
COPY catalog /maestro/catalog/

ENV LDAP_HOST ldap
ENV LDAP_PORT 389

EXPOSE 80

ENTRYPOINT ["/maestro/maestro"]
