package config

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/brickman-source/golang-utilities/viper"
	"github.com/fsnotify/fsnotify"
	"github.com/naoina/toml"

	"github.com/brickman-source/golang-utilities/converter"
	"github.com/brickman-source/golang-utilities/crypto"
	"github.com/brickman-source/golang-utilities/json"
	"github.com/brickman-source/golang-utilities/log"
	"github.com/brickman-source/golang-utilities/map/helper"
	gusync "github.com/brickman-source/golang-utilities/sync"

	_ "github.com/brickman-source/golang-utilities/viper/remote"
)

var (
	AesKeySalt string
)

type Config struct {
	cache  sync.Map
	v      *viper.Viper
	AesKey []byte
}

func NewConfig(path string) *Config {
	ret := &Config{
		v: viper.New(),
	}
	ret.v.SetConfigFile(path)
	err := ret.v.ReadInConfig()
	if err != nil {
		log.Errorf("read config file err:%s", err.Error())
	}
	ret.v.WatchConfig()
	ret.v.OnConfigChange(func(e fsnotify.Event) {
		log.Debugf("reload config")
		gusync.EraseSyncMap(&ret.cache)
	})
	return ret
}

func NewConfigWithString(configStr string) *Config {
	ret := &Config{
		v: viper.New(),
	}
	ret.v.SetConfigType("toml")
	err := ret.v.ReadConfig(strings.NewReader(configStr))
	if err != nil {
		log.Errorf("read config file err:%s", err.Error())
	}
	return ret
}

func ParseConfig(configStr string, configType string) *Config {
	ret := &Config{
		v: viper.New(),
	}
	ret.v.SetConfigType(configType)
	err := ret.v.ReadConfig(bytes.NewReader([]byte(configStr)))
	if err != nil {
		log.Errorf("read config err:%s", err.Error())
	}
	return ret
}

func NewRemoteConfig(provider, addr, path, configType string) *Config {
	ret := &Config{
		v: viper.New(),
	}
	ret.v.SetConfigType(configType)
	err := ret.v.AddRemoteProvider(provider, addr, path)
	if err != nil {
		log.Errorf("read log file err:%s", err.Error())
		return nil
	}
	err = ret.v.ReadRemoteConfig()
	if err != nil {
		log.Errorf("read log file err:%s", err.Error())
		return nil
	}
	err = ret.v.WatchRemoteConfigOnChannel()
	if err != nil {
		log.Infof("WatchRemoteConfigOnChannel err: %s", err.Error())
		return nil
	}
	ret.v.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("reload config: %s", string(json.ShouldMarshal(ret.v.AllSettings())))
		gusync.EraseSyncMap(&ret.cache)
	})

	return ret
}

func (cfg *Config) BindEnv(input ...string) error {
	return cfg.v.BindEnv(input...)
}

func (cfg *Config) SetDefault(key string, value interface{}) {
	cfg.v.SetDefault(key, value)
}

func (cfg *Config) GetString(path string) string {
	if v, ok := cfg.cache.Load(path); ok {
		if c, ok := v.(string); ok {
			return c
		}
	}

	ret := cfg.v.GetString(path)
	if strings.HasPrefix(ret, "aes://") {
		ret = cfg.DecryptString(ret[6:])
	}
	cfg.cache.Store(path, ret)
	//[path] = ret
	return ret
}

func (cfg *Config) ExtractStringFromMap(m map[string]interface{}, path string, defaultValue string) string {
	ret := helper.GetValueAsString(m, path, "")
	if strings.HasPrefix(ret, "aes://") {
		ret = cfg.DecryptString(ret[6:])
	}
	return ret
}

func (cfg *Config) GetInt(path string) int {
	if v, ok := cfg.cache.Load(path); ok {
		if c, ok := v.(int); ok {
			return c
		}
	}
	ret := cfg.v.GetInt(path)
	cfg.cache.Store(path, ret)
	return ret
}

func (cfg *Config) GetInt32(path string) int32 {
	if v, ok := cfg.cache.Load(path); ok {
		if c, ok := v.(int32); ok {
			return c
		}
	}
	ret := cfg.v.GetInt32(path)
	cfg.cache.Store(path, ret)
	return ret
}

func (cfg *Config) GetInt64(path string) int64 {
	if v, ok := cfg.cache.Load(path); ok {
		if c, ok := v.(int64); ok {
			return c
		}
	}
	ret := cfg.v.GetInt64(path)
	cfg.cache.Store(path, ret)
	return ret
}

func (cfg *Config) GetFloat64(path string) float64 {
	if v, ok := cfg.cache.Load(path); ok {
		if c, ok := v.(float64); ok {
			return c
		}
	}
	ret := cfg.v.GetFloat64(path)
	cfg.cache.Store(path, ret)
	return ret
}

func (cfg *Config) GetBool(path string) bool {
	if v, ok := cfg.cache.Load(path); ok {
		if c, ok := v.(bool); ok {
			return c
		}
	}
	ret := cfg.v.GetBool(path)
	cfg.cache.Store(path, ret)
	return ret
}

func (cfg *Config) GetDuration(path string) time.Duration {
	if v, ok := cfg.cache.Load(path); ok {
		if c, ok := v.(time.Duration); ok {
			return c
		}
	}
	ret := cfg.v.GetDuration(path)
	cfg.cache.Store(path, ret)
	return ret
}

func (cfg *Config) GetStringSlice(path string) []string {
	if v, ok := cfg.cache.Load(path); ok {
		if r, ok := v.([]string); ok {
			return r
		}
	}

	ret := cfg.v.GetStringSlice(path) //map_helper.GetValue(cfg.data, path)
	for i, item := range ret {
		if strings.HasPrefix(item, "aes://") {
			ret[i] = cfg.DecryptString(item[6:])
		}
	}
	cfg.cache.Store(path, ret)
	return ret
}

func (cfg *Config) GetStringMapString(path string) map[string]string {
	if v, ok := cfg.cache.Load(path); ok {
		if r, ok := v.(map[string]string); ok {
			return r
		}
	}
	ret := cfg.v.GetStringMapString(path) //map_helper.GetValue(cfg.data, path)
	cfg.cache.Store(path, ret)
	return ret
}

func (cfg *Config) GetMapSlice(path string) []map[string]interface{} {
	if v, ok := cfg.cache.Load(path); ok {
		if r, ok := v.([]map[string]interface{}); ok {
			return r
		}
	}
	mapSlice := cfg.v.Get(path)
	ret := converter.AsMapSlice(mapSlice)

	//ret, ok := cfg.v.Get(path).([]map[string]interface{})
	//if !ok {
	//	return nil
	//}

	cfg.cache.Store(path, ret)
	return ret
}

func (cfg *Config) GetValue(path string) interface{} {
	if v, ok := cfg.cache.Load(path); ok {
		return v
	}
	ret := cfg.v.Get(path)
	cfg.cache.Store(path, ret)
	return ret
}

func (cfg *Config) SetValue(path string, val interface{}) {
	cfg.v.Set(path, val)
	cfg.cache.Store(path, val)
}

func (cfg *Config) DecryptString(str string) string {
	if cfg.AesKey == nil {
		aesKey1 := cfg.GetString("crypto.aesKey")
		cfg.AesKey = crypto.String(aesKey1 + AesKeySalt).GetMd5()
	}

	eData, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Errorf("config DecryptString(%s) err:%s", str, err.Error())
		return ""
	}
	dData, err := crypto.AESDecryptCBC(eData, cfg.AesKey, cfg.AesKey[:aes.BlockSize])
	if err != nil {
		log.Errorf("config DecryptString(%s) err:%s", str, err.Error())
		return ""
	}

	return string(dData)
}

func (cfg *Config) EncryptString(str string) string {
	if cfg.AesKey == nil {
		aesKey1 := cfg.GetString("crypto.aesKey")
		cfg.AesKey = crypto.String(aesKey1 + AesKeySalt).GetMd5()
	}

	dData := []byte(str)
	eData, err := crypto.AESEncryptCBC(dData, cfg.AesKey, cfg.AesKey[:aes.BlockSize])
	if err != nil {
		log.Fatalf("config DecryptString(%v) err:%v", str, err)
	}

	return base64.StdEncoding.EncodeToString(eData)
}

func (cfg *Config) GetAllKeys() []string {
	return cfg.v.AllKeys()
}

func (cfg *Config) Save(path string) error {
	data, err := toml.Marshal(cfg.v.AllSettings())
	if err != nil {
		return err
	}
	if path != "" {
		err = ioutil.WriteFile(path, data, 0666)
	}
	return err
}

func (cfg *Config) Exit() {
	cfg.v.Exit()
}
