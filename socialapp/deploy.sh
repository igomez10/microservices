cd microservices/socialapp
git pull
docker compose build
docker compose down
docker compose up -d
docker builder prune -f
/usr/local/go/bin/go clean -modcache 
