package models

type PlayerFunctions struct {
	PlayPause                func() (bool, error)
	GetAllPlayers            func() ([]string, error)
	SeekForward              func() (bool, error)
	SeekBackward             func() (bool, error)
	SeekTo                   func(position int64) (bool, error)
	Next                     func() (bool, error)
	Prev                     func() (bool, error)
	SetVolume                func(volume int) (bool, error)
	GetCurrentPlayerMetadata func() (MediaInfo, error)
}

type MediaInfo struct {
	Title    string
	Artist   string
	Album    string
	Artwork  string
	Length   string
	TrackId  string
	Status   string
	Position string
	Player   string
}
