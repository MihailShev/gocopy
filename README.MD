# Домашнее задание №9
## Реализовать утилиту копирования файлов  (см man dd)

Выводить в консоль прогресс копирования. Программа
должна корректно обрабатывать ситуацию, когда offset или offset+limit за пределами source файла.

``` 
# копирует 2К из source в dest, пропуская 1K данных
$  gocopy ­from /path/to/source ­to /path/to/dest ­offset 1024 ­limit 2048
```
Настроить и запустить линтеры, создать Makele для автоматизации тестирования и сборки. Должна быть
возможность скачать протестировать и установить программу с помощью go get/test/install