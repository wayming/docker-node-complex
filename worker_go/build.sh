export REDIS_HOST=172.19.0.2
export REDIS_PORT=6379

go build -o bin/fibcalc src/fibcalc/main.go
./bin/fibcalc
