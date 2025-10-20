# golangci-lint v2.4.0

## Быстрый старт
1. Установить golangci-lint v2.4.0
```shell
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.4.0
```
2. Проверить установку
```
golangci-lint --version
```
3. Запустить из корневого каталога
```
golangci-lint run ./..
```
4. Ожидаемый результат
```
...
util/util.go:237:3: QF1003: could use tagged switch on caseName (staticcheck)
                if caseName == Nominative {
                ^
util/util.go:243:3: QF1003: could use tagged switch on caseName (staticcheck)
                if caseName == Nominative {
                ^
util/util.go:249:3: QF1003: could use tagged switch on caseName (staticcheck)
                if caseName == Nominative {
                ^
15 issues:
* errorlint: 1
* gocritic: 2
* staticcheck: 12
```

# Запуск через docker
Запускать лучше на конкретном пакете как в примерах команд запуска, 
**весь проект целиком сканирует плохо**, причина неизвестна. 

Локально должна быть создана папка **golangci-lint** (она добавлена в .gitignore и пушиться не будет).
В ней будут храниться кэщ и отчеты.

Можно добавить флаг **-v** в команду запуска, тогда будут выводиться логи golangci-lint.

## Запуск под Linux  

Запускается из корня проекта

```shell
docker run --rm \
    -v $(pwd):/client \
    -v $(pwd)/golangci-lint/.cache/golangci-lint/v2.4.0:/root/.cache \
    -w /client \
    golangci/golangci-lint:v2.4.0 \
        golangci-lint run ./... \
            -c .golangci.yml 
```
---
## Запуск под windows

**Не забыть запустить docker desktop.**

Запускается из корня проекта через **PowerShell**

```shell
docker run --rm `
    -v ${pwd}:/client `
    -v ${pwd}/golangci-lint/.cache/golangci-lint/v2.4.0:/root/.cache `
    -w /client `
    golangci/golangci-lint:v2.4.0 `
        golangci-lint run ./... `
            -c .golangci.yml 
```
---
## Запуск под Windows через WSL

Запускается из терминала wsl так же как и под линукс. 
  
### Если не запускается:


* Если повезёт, то хватит [этой статьи](https://docs.docker.com/desktop/features/wsl/)
* Если не повезло - придется накатить docker отдельно на WSL

---

## В случае проблем с линтерами

1. Попробовать обновить образ: `docker pull golangci/golangci-lint:v2.4.0` (LINUX/WSL)