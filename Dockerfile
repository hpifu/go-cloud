FROM centos:centos7
COPY docker/cloud /var/docker/cloud
RUN mkdir -p /var/docker/cloud/log
EXPOSE 6061
WORKDIR /var/docker/cloud
CMD [ "bin/cloud", "-c", "configs/cloud.json" ]
