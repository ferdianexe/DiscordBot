package music

import "sync"

type PlaylistStatus struct {
	PlayMutex     sync.RWMutex
	PlayListMutex sync.RWMutex
	SkipMutex     sync.RWMutex
	IsPlaying     bool
	IsSkipped     bool
	PlaylistURL   []string
}
