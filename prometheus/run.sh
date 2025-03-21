docker rm -f prom-exporter prometheus grafana

docker run -d -p 9100:9100 \
 --name=prom-exporter \
 --privileged=true \
 -v "/home/ken/ken/prometheus/host/proc:/host/proc:ro" \
 -v "/home/ken/ken/prometheus/host/sys:/host/sys:ro" \
 -v "/home/ken/ken/prometheus/rootfs:/rootfs:ro" \
 prom/node-exporter

 docker run  -d \
 --name=prometheus \
 --privileged=true \
 -p 9090:9090 \
 -v /home/ken/ken/prometheus/etc/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml  \
 prom/prometheus

docker run -d \
 -p 3000:3000 \
 --privileged=true \
 --name=grafana \
 -v /home/ken/ken/prometheus/var/lib/grafana-storage:/var/lib/grafana \
 grafana/grafana