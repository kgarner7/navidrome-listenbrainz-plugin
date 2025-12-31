//go:build wasip1

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/extism/go-pdk"
	"github.com/navidrome/navidrome/plugins/pdk/go/metadata"
)

const (
	lbzEndpoint = "https://api.listenbrainz.org/1/"

	labsBase  = "https://labs.api.listenbrainz.org/"
	algorithm = "session_based_days_9000_session_300_contribution_5_threshold_15_limit_50_skip_30"
)

var notFound = errors.New("not found")

type ListenBrainzAgent struct{}

// Ensure wikimediaPlugin implements the provider interfaces
var (
	_ metadata.ArtistURLProvider      = (*ListenBrainzAgent)(nil)
	_ metadata.ArtistTopSongsProvider = (*ListenBrainzAgent)(nil)
	_ metadata.SimilarArtistsProvider = (*ListenBrainzAgent)(nil)
)

func listenBrainzRequest(endpoint string, params url.Values) ([]byte, error) {
	url := fmt.Sprintf("%s%s?%s", lbzEndpoint, endpoint, params.Encode())
	req := pdk.NewHTTPRequest(pdk.MethodGet, url)
	req.SetHeader("Accept", "application/json")
	req.SetHeader("User-Agent", "NavidromeListenBrainzPlugin/0.1")

	resp := req.Send()

	if resp.Status() != 200 {
		return nil, fmt.Errorf("ListenBrainz HTTP error: status %d, body: %s", resp.Status, string(resp.Body()))
	}

	return resp.Body(), nil
}

type trackInfo struct {
	RecordingName string `json:"recording_name"`
	RecordingMbid string `json:"recording_mbid"`
}

func (l ListenBrainzAgent) GetArtistTopSongs(req metadata.TopSongsRequest) (metadata.TopSongsResponse, error) {
	if req.MBID == "" {
		return metadata.TopSongsResponse{}, notFound
	}

	resp, err := listenBrainzRequest("popularity/top-recordings-for-artist/"+req.MBID, url.Values{})
	if err != nil {
		return metadata.TopSongsResponse{}, err
	}

	var tracks []trackInfo
	if err := json.Unmarshal(resp, &tracks); err != nil {
		return metadata.TopSongsResponse{}, fmt.Errorf("failed to parse ListenBrainz response: %w", err)
	}

	// Make sure we do not exceed the number of requested songs.
	count := min(len(tracks), int(req.Count))

	songs := make([]metadata.SongRef, count)
	for idx, track := range tracks[:count] {
		songs[idx] = metadata.SongRef{
			MBID: track.RecordingMbid,
			Name: track.RecordingName,
		}
	}

	return metadata.TopSongsResponse{Songs: songs}, nil
}

type artistMetadataResult struct {
	Rels struct {
		OfficialHomepage string `json:"official homepage,omitempty"`
	} `json:"rels,omitzero"`
}

func (l ListenBrainzAgent) GetArtistURL(req metadata.ArtistRequest) (metadata.ArtistURLResponse, error) {
	if req.MBID == "" {
		return metadata.ArtistURLResponse{}, notFound
	}

	params := url.Values{}
	params.Add("artist_mbids", req.MBID)

	resp, err := listenBrainzRequest("metadata/artist", params)
	if err != nil {
		return metadata.ArtistURLResponse{}, err
	}

	var result []artistMetadataResult
	if err := json.Unmarshal(resp, &result); err != nil {
		return metadata.ArtistURLResponse{}, fmt.Errorf("failed to parse ListenBrainz response: %w", err)
	}

	if len(result) != 1 {
		return metadata.ArtistURLResponse{}, notFound
	}

	if result[0].Rels.OfficialHomepage != "" {
		return metadata.ArtistURLResponse{URL: result[0].Rels.OfficialHomepage}, nil
	}

	return metadata.ArtistURLResponse{}, notFound
}

type artist struct {
	MBID string `json:"artist_mbid"`
	Name string `json:"name"`
}

func (l ListenBrainzAgent) GetSimilarArtists(req metadata.SimilarArtistsRequest) (metadata.SimilarArtistsResponse, error) {
	if req.MBID == "" {
		return metadata.SimilarArtistsResponse{}, notFound
	}

	url := fmt.Sprintf("%ssimilar-artists/json?artist_mbids=%s&algorithm=%s", labsBase, req.MBID, algorithm)
	httpReq := pdk.NewHTTPRequest(pdk.MethodGet, url)
	httpReq.SetHeader("Accept", "application/json")
	httpReq.SetHeader("User-Agent", "NavidromeListenBrainzPlugin/0.1")

	resp := httpReq.Send()

	if resp.Status() != 200 {
		return metadata.SimilarArtistsResponse{}, fmt.Errorf("ListenBrainz labs HTTP error: status %d, body: %s", resp.Status, string(resp.Body()))
	}

	var lbzArtists []artist
	if err := json.Unmarshal(resp.Body(), &lbzArtists); err != nil {
		return metadata.SimilarArtistsResponse{}, fmt.Errorf("failed to parse ListenBrainz response: %w", err)
	}

	// Make sure we do not exceed the number of requested songs.
	count := min(len(lbzArtists), int(req.Limit))

	artists := make([]metadata.ArtistRef, count)
	for i, artist := range lbzArtists[:count] {
		artists[i] = metadata.ArtistRef{
			MBID: artist.MBID,
			Name: artist.Name,
		}
	}

	return metadata.SimilarArtistsResponse{
		Artists: artists,
	}, nil
}

func main() {}

func init() {
	metadata.Register(&ListenBrainzAgent{})
}
