Start:
```
docker run -d --name some-clickhouse-server --ulimit nofile=262144:262144 yandex/clickhouse-server
docker run -it --rm --link some-clickhouse-server:clickhouse-server yandex/clickhouse-client --host clickhouse-server
```

Stop and clean:
```
docker rm -f -v $(docker ps | grep some-clickhouse-server | awk '{ print $1}')
```
