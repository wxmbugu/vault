# vault

This is a share secret api

## Prerequisite
- go
- docker
- make
- migrate

## To get started:
```
git clone https://github.com/Wambug/vault
``` 
build :
```
go build .
```
setp postgres:
```
make postgres
```
creatdb:
```
make creatdb
```
migrations:
```
make migrateup
```
run server:
```
./secret-vault
```

client setup:
```
git clone https://github.com/Wambug/vaul-client
mv -r vaul-client /path/to/templates/vaul-client
```

start client:
```
make client
```
