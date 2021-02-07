/*
 * Copyright (c) 2018. Brickman Source.
 */

package baidu

import (
	"github.com/brickman-source/golang-utilities/cache"
	"github.com/brickman-source/golang-utilities/config"
	"sync"
)

type Baidu struct {
	cache  cache.Cache
	memory *sync.Map
	config *config.Config
	quit   chan struct{}
}

func NewBaidu(cah cache.Cache, config *config.Config) *Baidu {
	ret := &Baidu{
		cache:  cah,
		memory: new(sync.Map),
		config: config,
		quit:   make(chan struct{}),
	}

	return ret
}

func (bd *Baidu) Exit() error {
	close(bd.quit)
	return nil
}
