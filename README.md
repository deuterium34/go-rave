# go-rave

![пример в повер шелл](/preview.png)

Это простая TUI утилита для прослушивания музыки.

Установка `go install github.com/deuterium34/go-rave`

Поддерживаемые форматы:

| формат | состояние |
| ------ |:---------:|
|  .mp3  | ✅   |
|  .flac | TODO |
|  .ogg  | TODO |
|  .opus | TODO |
|  .wma  | TODO |
|  .wav  | TODO |
|  .mp4a | TODO |

## Использование

```
./go-rave.exe ./music.mp3
```

Для проигрывания в моно:

```
-m, --mono
```

Можно использовать семпрейт 44100 либо 48000, по умолчанию используется 44100:

```
-s, --sample-rate INT
```

А так-же:

```
-h, --help
--version
```
