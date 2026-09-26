# asiatorrents-fstop

[filesharetop](https://github.com/RangelReale/filesharetop) plugin for AsiaTorrents (asiatorrents.me). It scrapes the tracker's top torrent lists and serves a ranked "AsiaTorrents Top" website.

## Components

- **`asiatorrents-importer`** scrapes the tracker and stores a snapshot in MongoDB (database `fstop_asiatorrents`). It then recomputes the 48-hour and weekly rankings. Run it periodically, for example hourly.
- **`asiatorrents-site`** is the web UI, served on port **13111**.

Both binaries connect to MongoDB on the local host and accept `-version`.

## Configuration

AsiaTorrents requires a login. Copy `asiatorrents.conf` to `asiatorrents-current.conf` (this copy is gitignored) and fill in the cookies from a logged-in browser: `PHPSESSID`, `lastseen`, `pass`, `uid`. Pass it with `-configfile asiatorrents-current.conf`.

## Build

This is pre-modules Go code. Put this repo and `filesharetop` in a GOPATH workspace, then build:

```sh
go get github.com/RangelReale/asiatorrents-fstop/...
# or, from a GOPATH checkout:
go build ./asiatorrents-importer ./asiatorrents-site
```

## Run

```sh
# once per hour, e.g. from cron
./asiatorrents-importer/asiatorrents-importer -configfile asiatorrents-current.conf

# web UI at http://localhost:13111/
./asiatorrents-site/asiatorrents-site
```

## Tests

```sh
go test -run TestFetcher .
```

The test scrapes the live tracker, so it needs network access. It currently doesn't compile: `at_test.go` calls `NewFetcher()` without the required config.
