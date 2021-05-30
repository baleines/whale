# Whale

[![Build Status](https://github.com/pmalhaire/whale/workflows/run%20tests/badge.svg)](https://github.com/pmalhaire/whale/actions?workflow=run%20tests)
[![codecov](https://codecov.io/gh/codecov/example-go/branch/main/graph/badge.svg)](https://codecov.io/gh/codecov/example-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/pmalhaire/whale)](https://goreportcard.com/report/github.com/pmalhaire/whale)

Whale is a fun simple card game for 2 to 4 players.



## curent status

The game is in development and works with small limitations :

Bonuses are not implemented yet.

## gif overview

![whaleGame](whale.gif)

## build

```
go build -o whaleGame main.go
```

## run

from binary

```
./whaleGame
```

from source

```
go run main.go
```

## Roadmap

- fix asciinema to suport utf8
- implement bonuses
- split package to make it a lib
- implement a web version using view-js
