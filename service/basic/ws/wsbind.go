package ws

import "sync"

type binder struct {
	mutex     sync.RWMutex
	userIDMap map[string]string
}

func (b *binder) Bind(userAccount string) error {
	if userAccount == "" {

	}
	return nil
}
