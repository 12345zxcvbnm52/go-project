
 docker run -p 3307:3306 --name mysql \
 --privileged=true \
 -v /home/ken/ken/mysql/log:/var/log/mysql \
 -v /home/ken/ken/mysql/data:/var/lib/mysql \
 -v /home/ken/ken/mysql/conf/my.cnf:/etc/mysql/my.cnf \
 -v /home/ken/ken/mysql/mysql-files:/var/lib/mysql-files \
 -e MYSQL_ROOT_PASSWORD=123 -d mysql:8.0.20

 