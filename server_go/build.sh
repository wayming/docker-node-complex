export REDIS_HOST=172.19.0.7
export REDIS_PORT=6379
export PGUSER=postgres
export PGHOST=172.19.0.4
export PGDATABASE=postgres
export PGPASSWORD=postgres_password
export PGPORT=5432
go build -o bin/httpserver src/httpserver/main.go
./bin/httpserver