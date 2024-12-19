 docker run -itd --name transfer --privileged=true -p 8060:8060 \
  -v /home/ken/ken/transfer/web/static:/app/web/static \
  -v /home/ken/ken/transfer/transfer:/app/transfer \
  -v /home/ken/ken/transfer/app.yml:/app/app.yml \
  go-mysql-transfer
 