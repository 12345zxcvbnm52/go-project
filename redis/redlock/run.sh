docker run -d \
 -p 17715:17715 \
 -p 27715:27715 \
 --name redlock1 \
 --privileged=true \
 -v /home/ken/ken/redis/redlock/redis17715.conf:/etc/redis/redis.conf \
 -v /home/ken/ken/redis/redlock:/redlock \
 redis:7.0.15 \
 redis-server /etc/redis/redis.conf

 docker run -d \
 -p 17716:17716 \
 -p 27716:27716 \
 --name redlock2 \
 --privileged=true \
 -v /home/ken/ken/redis/redlock/redis17716.conf:/etc/redis/redis.conf \
 -v /home/ken/ken/redis/redlock:/redlock \
 redis:7.0.15 \
 redis-server /etc/redis/redis.conf

 docker run -d \
 -p 17717:17717 \
 -p 27717:27717 \
 --name redlock3 \
 --privileged=true \
 -v /home/ken/ken/redis/redlock/redis17717.conf:/etc/redis/redis.conf \
 -v /home/ken/ken/redis/redlock:/redlock \
 redis:7.0.15 \
 redis-server /etc/redis/redis.conf

  docker run -d \
 -p 17719:6379 \
 -p 27719:16379 \
 --name redlock4 \
 --privileged=true \
 -v /home/ken/ken/redis/redlock/redis17719.conf:/etc/redis/redis.conf \
 -v /home/ken/ken/redis/redlock:/redlock \
 redis:7.0.15 \
 redis-server /etc/redis/redis.conf

   docker run -d \
 -p 17718:6379 \
  -p 27718:16379 \
 --name redlock5 \
 --privileged=true \
 -v /home/ken/ken/redis/redlock/redis17718.conf:/etc/redis/redis.conf \
 -v /home/ken/ken/redis/redlock:/redlock \
 redis:7.0.15 \
 redis-server /etc/redis/redis.conf

   docker run -d \
 -p 17720:6379 \
 -p 27720:16379 \
 --name redlock6 \
 --privileged=true \
 -v /home/ken/ken/redis/redlock/redis17720.conf:/etc/redis/redis.conf \
 -v /home/ken/ken/redis/redlock:/redlock \
 redis:7.0.15 \
 redis-server /etc/redis/redis.conf
# redis-cli -a 123 \ 
#  --cluster create \
#  --cluster-replicas 1 \
#  192.168.199.128:17715 \
#  192.168.199.128:17716 \
#  192.168.199.128:17717 \
#  192.168.199.128:17719
