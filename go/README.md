# FastCat (Go edition)

Порт Bash-скриптов FastCat на Go: **один статический бинарник вместо
`bash` + `curl` + `git` + `sed` + `unzip`** (нужен только сам `fastfetch`).

## Сборка

```bash
cd go
go build -o fastcat .
```

Кросс-компиляция (сборка на Linux под любую ОС/архитектуру):

```bash
GOOS=linux  GOARCH=amd64 go build -o fastcat-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o fastcat-macos-arm64 .
GOOS=windows GOARCH=amd64 go build -o fastcat.exe .
```

Важно: бинарник должен лежать в корне репозитория FastCat (или запускаться из него),
так как темы (`Small-Themes/`, `Large-Themes/`, `Visuals-Themes/`) ищутся
относительно текущей директории и расположения бинарника (с подъёмом вверх по дереву).

## Соответствие Bash → Go

| Bash-скрипт | Замена |
|---|---|
| `fastcat.sh` | `go build → ./fastcat --shell` |
| `backup.sh` | `./fastcat --backup` (или пункт `[B]` в меню) |
| `install-icons.sh` | `./fastcat --icons` (распаковка zip — чистый Go, `unzip` не нужен) |
| `*/start.sh` | встроено в меню (`--shell`) |
| проверка версии через `curl` | встроенный HTTP-клиент, таймаут 3 c, не блокирует запуск |
| `sed`-правки конфигов | `strings`/`regexp` в `transformVisualConfig`, `rewritePreviewConfig` |

## Команды

```
--shell    интерактивное меню (Small / Large / Visuals)
--backup   бэкап ~/.config/fastfetch в ./Backup-<дата>
--icons    установка CascadiaCode Nerd Font
--update   синхронизация репозитория через git + подсказка пересобрать бинарник
--version  версия
--about    о проекте
--help     справка
```

## Отличия от Bash-версии

- **Списки тем не хардкодятся** — директории сканируются автоматически,
  новые темы появляются в меню без правок кода.
- В меню добавлены `[G]Gallery` — интерактивная галерея тем: листайте превью стрелками ←/→ (или j/k), Enter применяет, R — случайная; `[S]Search` (поиск тем), `[I]Icons`, `[U]Update`.
- Проверка обновлений сравнивает с `go/VERSION` на GitHub (при запуске — молча,
  с таймаутом 3 с), а не с версией bash-скрипта.
  **После пуша этой папки в репозиторий** не забудьте обновлять файл `go/VERSION`.
- Превью тем работает как в оригинале: конфиг копируется во временную папку,
  пути к ассетам перезаписываются.
- `read -n1` («нажмите любую клавишу») реализован через raw-режим терминала
  на Linux/macOS; на прочих ОС — обычный Enter.
- Обновление репозитория через `--update` синхронизирует **исходники**;
  бинарник нужно пересобрать (`go build -o fastcat .`).

## Структура

```
main.go    — вход, CLI-аргументы, главное меню, help/about
ui.go      — ANSI-цвета, баннеры, сетка тем, ввод, лоадер
themes.go  — сканирование тем, применение/превью/поиск/рандом, Visuals-протоколы
util.go    — копирование файлов, поиск корня репозитория, fastfetch, бэкап
update.go  — проверка и синхронизация обновлений (HTTP + git)
icons.go   — скачивание и установка Nerd-шрифта (HTTP + archive/zip)
key_*.go   — чтение одиночной клавиши (Linux/macOS/фолбэк)
```
