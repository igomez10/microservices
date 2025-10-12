set -e

cd microservices/socialapp
echo "Deploying socialapp..."
git pull
echo "Building and deploying socialapp..."
docker compose build
echo "Starting socialapp..."
docker compose up -d --remove-orphans
echo "Cleaning up..."
docker builder prune -f
echo "Cleaning Go module cache..."
/usr/local/go/bin/go clean -modcache 

cd ..

echo "Deploying urlshortener..."
cd urlshortener
echo "Pulling latest changes for urlshortener..."
git pull
echo "Building and deploying urlshortener..."
docker compose build
echo "Starting urlshortener..."
docker compose up -d --remove-orphans
echo "Cleaning up..."
docker builder prune -f
echo "Cleaning Go module cache..."
/usr/local/go/bin/go clean -modcache
