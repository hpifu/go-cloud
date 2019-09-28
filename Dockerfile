FROM centos:centos7
COPY docker/go-cloud /var/docker/go-cloud
RUN mkdir -p /var/docker/go-cloud/log
EXPOSE 6061
WORKDIR /var/docker/go-cloud
CMD [ "bin/cloud", "-c", "configs/cloud.json" ]
