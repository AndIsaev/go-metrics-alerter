# cmd/server

В данной директории будет содержаться код Сервера, который скомпилируется в бинарное приложение


### TRUSTED SUBNET
Для использования доверенной подсети используется флаг `-t` или переменная окружения `TRUSTED_SUBNET`
```shell
./server -t "127.0.0.1/32"
```
Важно! Сервер ожидает адрес в формате `CIDR` через `/` как предсавленно в примере, иначе клиент получит ошибку `403` 
