package wechat

import (
	"github.com/brickman-source/golang-utilities/cache"
	"github.com/brickman-source/golang-utilities/config"
)

type Wechat struct {
	cache  cache.Cache
	config *config.Config
	quit   chan struct{}
}

func NewWechat(cah cache.Cache, config *config.Config) *Wechat {
	ret := &Wechat{
		cache:  cah,
		config: config,
		quit:   make(chan struct{}),
	}
	go func() {
		ret.fetchAccessTokensLoop()
	}()
	return ret
}

func (wx *Wechat) Exit() error {
	close(wx.quit)
	return nil
}
