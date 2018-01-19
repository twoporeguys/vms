FROM redis:latest
MAINTAINER Pierre Chaussalet <p@2pg.bio>

ENV VMS_REDIS 127.0.0.1:6379
ENV VMS_PORT 8080

COPY bin/vms /usr/local/bin/vms
COPY vmsd.sh /usr/local/sbin/vmsd

EXPOSE 8080

CMD ["/usr/local/sbin/vmsd"]
