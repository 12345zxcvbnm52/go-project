docker run -d --name transfer \
 -v /home/ken/ken/transfer/app.yml:/app/app.yml \
 -p 8060:8060 transfer:latest 