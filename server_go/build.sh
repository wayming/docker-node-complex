export REDIS_HOST=192.168.16.2
export REDIS_PORT=6379
export PGUSER=postgres
export PGHOST=192.168.16.3
export PGDATABASE=postgres
export PGPASSWORD=postgres_password
export PGPORT=5432
go build -o bin/httpserver src/httpserver/main.go
./bin/httpserver
