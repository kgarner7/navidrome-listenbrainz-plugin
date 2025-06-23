# Navidrome ListenBrainz Plugin

A plugin for Navidrome for fetching the metadata from Listenbrainz:
- Artist homepage
- Artist top songs
- Similar artists

**NOTE**: Since the current implementation of `SimilarSongs` is serial, and each request to ListenBrainz is in series, this means that `getSimilarSongs` and `getSimilarSongs2` will take a while (upwards of 10 seconds)