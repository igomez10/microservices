cd microservices/socialapp
git pull
docker compose build
docker compose down
docker compose up -d
docker builder prune -f
go clean -modcache 
