FROM redis:7.0.15
MAINTAINER ken
ENV HOSTPATH /home/ken/ken/redis
WORKDIR /data
VOLUME $(HOSTPATH)/data:/data
#VOLUME $(HOSTPATH)/redis6381.conf:/etc/redis/redis6381.conf
COPY redis6381.conf /etc/redis/redis.conf
EXPOSE 6379
ENTRYPOINT ["redis-server"]
CMD ["/etc/redis/redis.conf"]
