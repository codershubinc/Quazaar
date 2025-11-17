package models

type SpotifyUser struct {
	Country         string `json:"country"`
	DisplayName     string `json:"display_name"`
	Email           string `json:"email"`
	ExplicitContent struct {
		FilterEnabled bool `json:"filter_enabled"`
		FilterLocked  bool `json:"filter_locked"`
	} `json:"explicit_content"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Href  string `json:"href"`
		Total int    `json:"total"`
	} `json:"followers"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Product string `json:"product"`
	Type    string `json:"type"`
	URI     string `json:"uri"`
}

type SpotifyFollowedArtists struct {
	Artists struct {
		Href    string `json:"href"`
		Limit   int    `json:"limit"`
		Next    string `json:"next"`
		Cursors struct {
			After  string `json:"after"`
			Before string `json:"before"`
		} `json:"cursors"`
		Total int `json:"total"`
		Items []struct {
			ExternalURLs map[string]string `json:"external_urls"`
			Followers    struct {
				Href  string `json:"href"`
				Total int    `json:"total"`
			} `json:"followers"`
			Genres []string `json:"genres"`
			Href   string   `json:"href"`
			ID     string   `json:"id"`
			Images []struct {
				URL    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name       string `json:"name"`
			Popularity int    `json:"popularity"`
			Type       string `json:"type"`
			URI        string `json:"uri"`
		} `json:"items"`
	} `json:"artists"`
}

type SpotifyTrack struct {
	Album struct {
		AlbumType        string            `json:"album_type"`
		TotalTracks      int               `json:"total_tracks"`
		AvailableMarkets []string          `json:"available_markets"`
		ExternalURLs     map[string]string `json:"external_urls"`
		Href             string            `json:"href"`
		ID               string            `json:"id"`
		Images           []struct {
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"images"`
		Name                 string `json:"name"`
		ReleaseDate          string `json:"release_date"`
		ReleaseDatePrecision string `json:"release_date_precision"`
		Type                 string `json:"type"`
		URI                  string `json:"uri"`
		Artists              []struct {
			ExternalURLs map[string]string `json:"external_urls"`
			Href         string            `json:"href"`
			ID           string            `json:"id"`
			Name         string            `json:"name"`
			Type         string            `json:"type"`
			URI          string            `json:"uri"`
		} `json:"artists"`
	} `json:"album"`
	Artists []struct {
		ExternalURLs map[string]string `json:"external_urls"`
		Href         string            `json:"href"`
		ID           string            `json:"id"`
		Name         string            `json:"name"`
		Type         string            `json:"type"`
		URI          string            `json:"uri"`
	} `json:"artists"`
	AvailableMarkets []string          `json:"available_markets"`
	DiscNumber       int               `json:"disc_number"`
	DurationMS       int               `json:"duration_ms"`
	Explicit         bool              `json:"explicit"`
	ExternalIDs      map[string]string `json:"external_ids"`
	ExternalURLs     map[string]string `json:"external_urls"`
	Href             string            `json:"href"`
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Popularity       int               `json:"popularity"`
	PreviewURL       string            `json:"preview_url"`
	TrackNumber      int               `json:"track_number"`
	Type             string            `json:"type"`
	URI              string            `json:"uri"`
	IsLocal          bool              `json:"is_local"`
}

type SpotifyPlaylist struct {
	Collaborative bool              `json:"collaborative"`
	Description   string            `json:"description"`
	ExternalURLs  map[string]string `json:"external_urls"`
	Href          string            `json:"href"`
	ID            string            `json:"id"`
	Images        []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Name  string `json:"name"`
	Owner struct {
		ExternalURLs map[string]string `json:"external_urls"`
		Href         string            `json:"href"`
		ID           string            `json:"id"`
		Type         string            `json:"type"`
		URI          string            `json:"uri"`
		DisplayName  string            `json:"display_name"`
	} `json:"owner"`
	Public     *bool  `json:"public"`
	SnapshotID string `json:"snapshot_id"`
	Tracks     struct {
		Href     string `json:"href"`
		Limit    int    `json:"limit"`
		Next     string `json:"next"`
		Offset   int    `json:"offset"`
		Previous string `json:"previous"`
		Total    int    `json:"total"`
		Items    []struct {
			AddedAt string `json:"added_at"`
			AddedBy struct {
				ExternalURLs map[string]string `json:"external_urls"`
				Href         string            `json:"href"`
				ID           string            `json:"id"`
				Type         string            `json:"type"`
				URI          string            `json:"uri"`
			} `json:"added_by"`
			IsLocal bool         `json:"is_local"`
			Track   SpotifyTrack `json:"track"`
		} `json:"items"`
	} `json:"tracks"`
	Type string `json:"type"`
	URI  string `json:"uri"`
}

type SpotifyPlaylistItems struct {
	Href     string `json:"href"`
	Limit    int    `json:"limit"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
	Items    []struct {
		AddedAt string `json:"added_at"`
		AddedBy struct {
			ExternalURLs map[string]string `json:"external_urls"`
			Href         string            `json:"href"`
			ID           string            `json:"id"`
			Type         string            `json:"type"`
			URI          string            `json:"uri"`
		} `json:"added_by"`
		IsLocal bool         `json:"is_local"`
		Track   SpotifyTrack `json:"track"`
	} `json:"items"`
}

type SpotifySavedTracks struct {
	Href     string `json:"href"`
	Limit    int    `json:"limit"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
	Items    []struct {
		AddedAt string       `json:"added_at"`
		Track   SpotifyTrack `json:"track"`
	} `json:"items"`
}

type SpotifyDevices struct {
	Devices []struct {
		ID               string `json:"id"`
		IsActive         bool   `json:"is_active"`
		IsPrivateSession bool   `json:"is_private_session"`
		IsRestricted     bool   `json:"is_restricted"`
		Name             string `json:"name"`
		Type             string `json:"type"`
		VolumePercent    int    `json:"volume_percent"`
		SupportsVolume   bool   `json:"supports_volume"`
	} `json:"devices"`
}
type SpotifyDeviceInfo struct {
	ID            string `json:"id"`
	IsActive      bool   `json:"is_active"`
	IsRestricted  bool   `json:"is_restricted"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	VolumePercent int    `json:"volume_percent"`
}

type SpotifyContext struct {
	ExternalURLs map[string]string `json:"external_urls"`
	Href         string            `json:"href"`
	Type         string            `json:"type"`
	URI          string            `json:"uri"`
}

type SpotifyPlaybackActions struct {
	Disallows map[string]bool `json:"disallows"`
}

type SpotifyCurrentlyPlaying struct {
	Device               SpotifyDeviceInfo      `json:"device"`
	RepeatState          string                 `json:"repeat_state"`
	ShuffleState         bool                   `json:"shuffle_state"`
	Context              SpotifyContext         `json:"context"`
	Timestamp            int64                  `json:"timestamp"`
	ProgressMS           int                    `json:"progress_ms"`
	IsPlaying            bool                   `json:"is_playing"`
	Item                 SpotifyTrack           `json:"item"`
	CurrentlyPlayingType string                 `json:"currently_playing_type"`
	Actions              SpotifyPlaybackActions `json:"actions"`
}
