psql -U postgres -h database -d postgres -c "CREATE DATABASE urlshortener;"
psql -U postgres -h database -d urlshortener -f /schema.sql
