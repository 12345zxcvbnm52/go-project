FROM redis:7.0.15
LABEL MAINTAINER="ken"
ENV HOSTPATH /home/ken/ken/redis
WORKDIR /data
#声明多个挂载点,但还需要在运行容器时指明挂载点,但是不能声明非文件夹的挂载点
VOLUME ["/data","/data/cluster"]

#更推荐COPY而非ADD,sredis默认是slave节点的配置文件
#默认使用普通的redis配置文件,如果要用集群就把redis.conf挂载data内的几个带标签的.conf
COPY redis.conf /etc/redis/redis.conf

# 创建一个新用户并切换到该用户,防止权限问题
#RUN useradd -ms /bin/bash myuser
#USER myuser
#仅仅是作为日志的形式告诉使用者这个dockerfile生成的镜像应当运行在哪一个端口上(容器内)
EXPOSE 6379
ENTRYPOINT ["redis-server"]
CMD ["/etc/redis/redis.conf"]
