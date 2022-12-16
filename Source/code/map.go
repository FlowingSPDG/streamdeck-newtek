package sdnewtek

import (
	"image"
	"sync"
)

// cacher Image cacher
type cacher struct {
	m sync.Map
}

type cacherKey struct {
	host string
	name string
}

func (c *cacher) getCache(key cacherKey) (image.Image, bool) {
	v, ok := c.m.Load(key)
	if !ok {
		return nil, ok
	}
	i, o := v.(image.Image)
	if !o {
		return nil, o
	}
	return i, ok
}

func (c *cacher) storeCache(key cacherKey, img image.Image) {
	c.m.Store(key, img)
}
