package music

import "sync"

type PlaylistStatus struct {
	PlayMutex     sync.RWMutex
	PlayListMutex sync.RWMutex
	IsPlaying     bool
	PlaylistURL   []string
}
