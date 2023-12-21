#init docker kafka and zookeeper
docker-compose up -d
# exec docker
docker exec -it ID bash
# create topic
kafka-topics --bootstrap-server localhost:29093 --topic report-topic --create  --partitions 3 --replication-factor 2