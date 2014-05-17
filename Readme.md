# Picturelife Uploader

Scans the supplied directories on the command line and hashes them and uploads them to Picturelife.

Performance can easily be improved at the moment (threading, batching signature checks, figuring out how to stream multipart uploads).

## Building

It's standard Go. Doesn't use anything outside the standard Go libraries. `go build .` will do nicely. `[goxc](https://github.com/laher/goxc)` works as well if you want to build for other platforms (you can cross compile to Linux, Mac, Windows, BSD, ARM, plan9 etc.).

## Usage

```
Usage of ./picturelife_uploader:
  -base_endpoint="https://api.picturelife.com/": API base endpoint location.
  -cache_dir="./": Path to where to store hash cache.
  -token="": API access token.
```

The hash cache stores files containing the hash and hashed timestamp of each file. If the mod_time of a file changes, the hash is recalculated.

The arguments are paths to walk looking for files.