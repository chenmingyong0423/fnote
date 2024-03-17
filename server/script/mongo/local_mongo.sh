set -e
docker-compose -f script/mongo/local_mongo_compose.yml down -v
docker-compose -f script/mongo/local_mongo_compose.yml up -d