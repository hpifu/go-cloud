FROM centos:centos7

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" >> /etc/timezone

COPY docker/go-cloud /var/docker/go-cloud
RUN mkdir -p /var/docker/go-cloud/log

EXPOSE 6061

WORKDIR /var/docker/go-cloud
CMD [ "bin/cloud", "-c", "configs/cloud.json" ]
