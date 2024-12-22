#多行模式
# docker run -p 9000:9000 -p 9090:9090 \
#      --net=host \
#      --name minio \
#      -e "MINIO_ACCESS_KEY=kensame" \
#      -e "MINIO_SECRET_KEY=kensame" \
#      -v /home/ken/ken/minio/data:/data \
#      -v /home/ken/ken/minio/config:/root/.minio \
#      minio/minio \ 
#      server /data --console-address ":9090" -address ":9000"

#单行模式
docker run -p 9000:9000 -p 9090:9090  \
 --net=host      \
 --name minio      \
 -d --restart=always      \
 -e "MINIO_ACCESS_KEY=kensame"      \
 -e "MINIO_SECRET_KEY=zxcvbnm52"      \
 -v /home/ken/ken/minio/data:/data      \
 -v /home/ken/ken/minio/config:/root/.minio      \
 minio/minio \
 server /data --console-address ":9090" -address ":9000"
