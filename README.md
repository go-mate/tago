# tago
set tag to git repo with golang.

# Installation

```bash
go install github.com/go-mate/tago/cmd/tago@latest
```

# Usage

## Show tags
```bash
tago
```

## Bump tag va.b.c to tag va.b.c+1 and push new tag. when tag new tag it will ask you to confirm.
```bash
tago bump
```

## Bump tag va.b.c to tag va.b.c+1 and push new tag. without ask you to confirm.
```bash
tago bump -b=100
```

## License

MIT License - See the `LICENSE` file for more details.

## Thank you

Give me stars. Thank you!!!

[![starring](https://starchart.cc/go-mate/tago.svg?variant=adaptive)](https://starchart.cc/go-mate/tago)
