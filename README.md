# URL Shortener
простой REST API сервис по сокращению ссылок.

сделано по следующему видео: https://www.youtube.com/watch?v=rCJvW2xgnk0

# Сборка
## Локально
```cmd
$ git clone https://github.com/gl0ckchan/url_shortener.git
$ cd url_shortener/
$ go build -o url_shortener_cmd cmd/url-shortener/main.go 
$ CONFIG_PATH=./config/local.yaml ./url_shortener_cmd
```
## В docker
```cmd
$ git clone https://github.com/gl0ckchan/url_shortener.git
$ cd url_shortener/
$ docker build -t shortener .
$ docker run -d -p 6969:6969 --name shortener-container shortener
```

# TODO
- [ ] добавить поддержку Redis.
- [x] добавить удаление URL.
- [ ] добавить админ панель (статистика, состояние, удаление).
- [ ] добавить gRPC-сервис для авторизации.
- [ ] сделать вывод логов цветным. 
- [x] добавить Dockerfile.
