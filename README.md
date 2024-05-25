# Windows Setup

Download Mingw-w64 https://www.mingw-w64.org/downloads/ and add to go environment:

```console
go env -w CGO_ENABLED=1
go env -w CC=C:\mingw64\bin\gcc.exe
```

## Playing The Game
Then build the game and play it with (the build might take a long time the first time, beacuse it compiles open-gl stuff - just trust the process and wait!)

```console
go build ./play_game/
./play_game.exe
```

## Running The Bots

Move Sequence Bot

```console
go run ./bot/
```
Neural Network Bot

```console
go run ./bot_nn/
```




