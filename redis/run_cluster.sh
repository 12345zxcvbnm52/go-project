 docker run -d \
 --name myredis \
 -p 6381:6379 \
 -v /home/ken/ken/redis/data:/data \
 -v /home/ken/ken/redis/data/cluster:/data/cluster \
 --privileged=true cludis:1.0 /etc/redis/redis.conf
 

redis-cli -a 123 \ 
 --cluster create \
 --cluster-replicas 1 \
 192.168.199.128:6379 \
 192.168.199.128:6380 \
 192.168.199.128:6381 \
 192.168.199.128:6382 \
 192.168.199.128:6383 \
 192.168.199.128:6384
