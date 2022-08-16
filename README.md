# vault

This is a share secret api

## To get started:
```
git clone https://github.com/Wambug/vault
``` 
build :
```
go build .
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
mv -r vaul-client /path/to/templates/
```

start client:
```
make client
```