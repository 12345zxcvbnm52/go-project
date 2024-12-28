 docker run --name es -p 9200:9200 -p 9300:9300 \
 -e "discovery.type=single-node" \
 -e ES_JAVA_OPTS="-Xms128m -Xmx256m" \
 --privileged=true \
 -v /home/ken/ken/elasticsearch/config/elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml \
 -v /home/ken/ken/elasticsearch/data:/usr/share/elasticsearch/data \
 -v /home/ken/ken/elasticsearch/plugins:/usr/share/elasticsearch/plugins \
 -d elasticsearch:7.10.1

 docker run -d --name kibana \
 -e ELASTICSEARCH_HOSTS="http://192.168.199.128:9200" \
 -p 5601:5601 \
 kibana:7.10.1