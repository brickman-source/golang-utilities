/*
 * Copyright (c) 2018. Brickman Source.
 */

package baidu

import (
	"fmt"
	"github.com/brickman-source/golang-utilities/cache"
	"github.com/brickman-source/golang-utilities/config"
	"github.com/brickman-source/golang-utilities/log"
	"sync"
)

type Baidu struct {
	cache   cache.Cache
	memory  *sync.Map
	config  *config.Config
	logFunc func(string)
	quit    chan struct{}
}

func NewBaidu(cah cache.Cache, config *config.Config, logFunc func(string)) *Baidu {
	ret := &Baidu{
		cache:   cah,
		memory:  new(sync.Map),
		config:  config,
		logFunc: logFunc,
		quit:    make(chan struct{}),
	}

	return ret
}

func (bd *Baidu) Exit() error {
	close(bd.quit)
	return nil
}

func (bd *Baidu) logf(format string, args ...interface{}) {
	if bd.logFunc != nil {
		bd.logFunc(fmt.Sprintf(format, args...))
	}
	log.Infof(format, args...)
}
