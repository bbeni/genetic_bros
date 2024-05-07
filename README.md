# Windows Setup

Download Mingw-w64 https://www.mingw-w64.org/downloads/ and add to go environment:

```console
go env -w CGO_ENABLED=1
go env -w CC=C:\mingw64\bin\gcc.exe
```
Then build the gui with

```console
go build ./gui/
```