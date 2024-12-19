docker run -d \
 -p 6379:6379 \
 --name redis \
 --privileged=true \
 -v /home/ken/ken/redis/redis.conf:/etc/redis/redis.conf \
 -v /home/ken/ken/redis/data:/data \
 redis:7.0.15 \
 redis-server /etc/redis/redis.conf