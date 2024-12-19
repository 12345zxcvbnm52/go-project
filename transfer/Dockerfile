
FROM ubuntu

WORKDIR /app
COPY app.yml /app/
COPY web/statics /app/
COPY transfer /app/
RUN ls -al
# 注意修改端口
EXPOSE 8060

ENTRYPOINT ["./transfer"]