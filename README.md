# Navidrome ListenBrainz Plugin

A plugin for Navidrome for fetching the metadata from Listenbrainz:
- Artist homepage
- Artist top songs
- Similar artists

**NOTE**: Since the current implementation of `SimilarSongs` is serial, and each request to ListenBrainz is in series, this means that `getSimilarSongs` and `getSimilarSongs2` will take a while (upwards of 10 seconds)

## Requirements
Navidrome >= 0.60.0. This reworked the plugin API. When upgrading to this version, you will need to install a new plugin.
For older versions of Navidrome, see https://github.com/kgarner7/navidrome-listenbrainz-plugin/releases/tag/v1.0.2

## Install instructions

### From GitHub Release

You can download the `brainz.ndp` from the latest release and then run `navidrome plugin install brainz.ndp`.
Make sure to run this command as your navidrome user.
This will unzip the package, and install it automatically in your plugin directory.

### From source

Requirements:
- `go` 1.25
- [`tinygo`](https://tinygo.org/) (recommended)

#### Build WASM plugin

##### Using TinyGo
```bash
tinygo build -o plugin.wasm -target wasip1 -buildmode=c-shared .
zip brainz.ndp plugin.wasm manifest.json
```

##### Using stock golang

```bash
go mod download
GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o plugin.wasm plugin.go
zip brainz.ndp plugin.wasm manifest.json
```

#### Package plugin

Copy the following files: `manifest.json` and `plugin.wasm`. 
Put them in a directory in your Navidrome `Plugins.Folder`.
Make sure that:
1. You have plugins enabled (`Plugins.Enabled = true`, `ND_PLUGINS_ENABLED = true`).
2. Your Navidrome user has read permissions in the plugin directory
