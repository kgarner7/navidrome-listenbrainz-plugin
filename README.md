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

Put the the `listenbrainz-metadata-provider.ndp ` file in your Navidrome `Plugins.Folder`.

Make sure that:
1. You have plugins enabled (`Plugins.Enabled = true`, `ND_PLUGINS_ENABLED = true`).
2. Your Navidrome user has read permissions in the plugin directory

As an admin user open the plugin page (profile icon > plugins) and enable the ` listenbrainz-metadata-provider` plugin.

### From source

Requirements:
- `go` 1.25
- [`tinygo`](https://tinygo.org/) (recommended)

#### Build WASM plugin

##### Using stock golang

```bash
make
```

This is a development build of the plugin. Compilation should be _extremely_ fast

##### Using TinyGo
```bash
make prod
```

This is the production version of the plugin.
Expect compilation to be slower, but the binary is also slower.

### Install

Copy the package `listenbrainz-metadata-provider.ndp` to your Navidrome plugin directory.
As an admin user open the plugin page (profile icon > plugins) and enable the `listenbrainz-metadata-provider` plugin.

Add the plugin name (`listenbrainz-metadata-provider`) to your [Agents](https://navidrome.org/docs/usage/configuration/options/#:~:text=Default%20Value-,agents,-ND_AGENTS).
