FROM scratch
COPY maestro /
COPY catalog/* /catalog/
ENV LDAP_HOST ldap
ENV LDAP_PORT 389
EXPOSE 80
ENTRYPOINT ["/maestro"]
