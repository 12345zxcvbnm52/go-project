docker run -d \
 -p 17715:17715 \
 --name redlock1 \
 --privileged=true \
 -v /home/ken/ken/redis/redlock/redis17715.conf:/etc/redis/redis.conf \
 -v /home/ken/ken/redis/redlock:/redlock \
 redis:7.0.15 \
 redis-server /etc/redis/redis.conf

 docker run -d \
 -p 17716:17716 \
 --name redlock2 \
 --privileged=true \
 -v /home/ken/ken/redis/redlock/redis17716.conf:/etc/redis/redis.conf \
 -v /home/ken/ken/redis/redlock:/redlock \
 redis:7.0.15 \
 redis-server /etc/redis/redis.conf

 docker run -d \
 -p 17717:17717 \
 --name redlock3 \
 --privileged=true \
 -v /home/ken/ken/redis/redlock/redis17717.conf:/etc/redis/redis.conf \
 -v /home/ken/ken/redis/redlock:/redlock \
 redis:7.0.15 \
 redis-server /etc/redis/redis.conf