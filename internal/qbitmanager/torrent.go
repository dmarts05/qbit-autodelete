package qbitmanager

// Holds torrent information
type Torrent struct {
	Hash  string `json:"hash"`
	Name  string `json:"name"`
	State string `json:"state"`
}
