//go:build wasip1

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/navidrome/navidrome/plugins/api"
	"github.com/navidrome/navidrome/plugins/host/http"
)

const (
	lbzEndpoint  = "https://api.listenbrainz.org/1/"
	lbzTimeoutMs = 5000

	labsBase    = "https://labs.api.listenbrainz.org/"
	labsTimeout = 10000
	algorithm   = "session_based_days_9000_session_300_contribution_5_threshold_15_limit_50_skip_30"
)

var (
	client = http.NewHttpService()
)

type ListenBrainzAgent struct{}

func listenBrainzRequest(ctx context.Context, endpoint string, params url.Values) ([]byte, error) {
	url := fmt.Sprintf("%s%s?%s", lbzEndpoint, endpoint, params.Encode())
	httpReq := &http.HttpRequest{
		Url: url,
		Headers: map[string]string{
			"Accept":     "application/json",
			"User-Agent": "NavidromeListenBrainzPlugin/0.1",
		},
		TimeoutMs: lbzTimeoutMs,
	}

	resp, err := client.Get(ctx, httpReq)
	if err != nil {
		return nil, fmt.Errorf("ListenBrainz request error: %w", err)
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("ListenBrainz HTTP error: status %d, body: %s", resp.Status, string(resp.Body))
	}

	return resp.Body, nil
}

func (l ListenBrainzAgent) GetAlbumImages(context.Context, *api.AlbumImagesRequest) (*api.AlbumImagesResponse, error) {
	return nil, api.ErrNotImplemented
}

func (l ListenBrainzAgent) GetAlbumInfo(context.Context, *api.AlbumInfoRequest) (*api.AlbumInfoResponse, error) {
	return nil, api.ErrNotImplemented
}

func (l ListenBrainzAgent) GetArtistBiography(context.Context, *api.ArtistBiographyRequest) (*api.ArtistBiographyResponse, error) {
	return nil, api.ErrNotImplemented
}

func (l ListenBrainzAgent) GetArtistImages(context.Context, *api.ArtistImageRequest) (*api.ArtistImageResponse, error) {
	return nil, api.ErrNotImplemented
}

func (l ListenBrainzAgent) GetArtistMBID(context.Context, *api.ArtistMBIDRequest) (*api.ArtistMBIDResponse, error) {
	return nil, api.ErrNotImplemented
}

type trackInfo struct {
	RecordingName string `json:"recording_name"`
	RecordingMbid string `json:"recording_mbid"`
}

func (l ListenBrainzAgent) GetArtistTopSongs(ctx context.Context, req *api.ArtistTopSongsRequest) (*api.ArtistTopSongsResponse, error) {
	if req.Mbid == "" {
		return nil, api.ErrNotFound
	}

	resp, err := listenBrainzRequest(ctx, "popularity/top-recordings-for-artist/"+req.Mbid, url.Values{})
	if err != nil {
		return nil, err
	}

	var tracks []trackInfo
	if err := json.Unmarshal(resp, &tracks); err != nil {
		return nil, fmt.Errorf("failed to parse ListenBrainz response: %w", err)
	}

	// Make sure we do not exceed the number of requested songs.
	count := min(len(tracks), int(req.Count))

	songs := make([]*api.Song, count)
	for idx, track := range tracks[:count] {
		songs[idx] = &api.Song{
			Mbid: track.RecordingMbid,
			Name: track.RecordingName,
		}
	}

	return &api.ArtistTopSongsResponse{Songs: songs}, nil
}

type artistMetadataResult struct {
	Rels struct {
		OfficialHomepage string `json:"official homepage,omitempty"`
	} `json:"rels,omitzero"`
}

func (l ListenBrainzAgent) GetArtistURL(ctx context.Context, req *api.ArtistURLRequest) (*api.ArtistURLResponse, error) {
	if req.Mbid == "" {
		return nil, api.ErrNotFound
	}

	params := url.Values{}
	params.Add("artist_mbids", req.Mbid)

	resp, err := listenBrainzRequest(ctx, "metadata/artist", params)
	if err != nil {
		return nil, err
	}

	var result []artistMetadataResult
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse ListenBrainz response: %w", err)
	}

	if len(result) != 1 {
		return nil, api.ErrNotFound
	}

	if result[0].Rels.OfficialHomepage != "" {
		return &api.ArtistURLResponse{Url: result[0].Rels.OfficialHomepage}, nil
	}

	return nil, api.ErrNotFound
}

type artist struct {
	MBID string `json:"artist_mbid"`
	Name string `json:"name"`
}

func (l ListenBrainzAgent) GetSimilarArtists(ctx context.Context, req *api.ArtistSimilarRequest) (*api.ArtistSimilarResponse, error) {
	if req.Mbid == "" {
		return nil, api.ErrNotFound
	}

	url := fmt.Sprintf("%ssimilar-artists/json?artist_mbids=%s&algorithm=%s", labsBase, req.Mbid, algorithm)
	httpReq := &http.HttpRequest{
		Url: url,
		Headers: map[string]string{
			"Accept":     "application/json",
			"User-Agent": "NavidromeListenBrainzPlugin/0.1",
		},
		TimeoutMs: labsTimeout,
	}

	resp, err := client.Get(ctx, httpReq)
	if err != nil {
		return nil, fmt.Errorf("ListenBrainz labs request error: %w", err)
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("ListenBrainz labs HTTP error: status %d, body: %s", resp.Status, string(resp.Body))
	}

	var lbzArtists []artist
	if err := json.Unmarshal(resp.Body, &lbzArtists); err != nil {
		return nil, fmt.Errorf("failed to parse ListenBrainz response: %w", err)
	}

	// Make sure we do not exceed the number of requested songs.
	count := min(len(lbzArtists), int(req.Limit))

	artists := make([]*api.Artist, count)
	for i, artist := range lbzArtists[:count] {
		artists[i] = &api.Artist{
			Mbid: artist.MBID,
			Name: artist.Name,
		}
	}

	return &api.ArtistSimilarResponse{
		Artists: artists,
	}, nil
}

func main() {}

func init() {
	api.RegisterMetadataAgent(ListenBrainzAgent{})
}
