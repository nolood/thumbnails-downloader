# Thumbnails downloader
# Запуск

## Thumbs - прокси-сервер для получения ссылки на скачивание превью с видеороликов youtube
```shell
git clone git@github.com:nolood/echelon-test.git
cd echelon-test/thumbs
go mod tidy
go run cmd/thumbs/main.go --config=./config/local.yml 
```

## Cli-thumbs - клиент для взаимодействия с сервисом
```shell
cd ../cli-thumbs
go mod tidy
mkdir images
go run cmd/thumbs/main.go --async [link1] [link2]
```
- Скачанные превью будут в папке images (!!! Важно, не забыть её создать)

# Примечание

Рефакторинг написанного не проводился. Грязь в клиенте для скачивания превью. Многие моменты могли бы быть лучше.