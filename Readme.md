# shortener
Сервис для создания укороченных ссылок.

- [ ] Добавить кастомную структуру для работы с ошибками, чтобы можно было передавать их между слоями и логгировать на верхнем уровне в слое handler. 

## Загрузка

* Загрузите последнюю версию бинарного файла со страницы релизов. Также можете воспользоваться командой для скачивания:
```
TODO
```
* Или загрузити образ Docker из репозитория и запустите его:
```
docker pull ghcr.io/mtchuikov/shortener:latest
```
* Или загрузите исходный код и выполните компиляцию самостоятельно
```
git clone https://github.com/mtchuikov/shortener
go build -ldflags="-s -w" -o ./build/shortener ./cmd/main.go
```

## Конфигурация

В настоящий момент поддерживается настройка поведения программы при помощи флагов командной строки. В следующих версиях будет добавлена конфигурации при помощи файла.

| Флаг        | Краткое обозначение | Значение по умолчанию   | Описание                                |
|-------------|---------------------|-------------------------|-----------------------------------------|
| --server-addr      | -a                  | 127.0.0.1:8080          | Адрес сервера                           |
| --base-url      | -b                  | http://127.0.0.1:8080   | Базовый URL для сокращённых ссылок      |
| --verbose   | -v                  | false                   | Вывод подробных логов                   |
