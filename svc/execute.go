package svc

import (
	"net/url"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"time"

	goSvc "github.com/judwhite/go-svc"

	"github.com/brickman-source/golang-utilities/config"
	"github.com/brickman-source/golang-utilities/log"
)

type Executor struct {
	s   Service
	cfg *config.Config
}

func Execute(s Service) {
	e := &Executor{}
	e.s = s
	if err := goSvc.Run(e, syscall.SIGINT, syscall.SIGTERM); err != nil {
		panic(err)
	}
}

func (e *Executor) Init(env goSvc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		err := os.Chdir(dir)
		if err != nil {
			panic(err.Error())
		}
	}
	remoteConfigURL := os.Getenv("REMOTE_CONFIG_URL")
	localConfigURL := os.Getenv("LOCAL_CONFIG_URL")

	if remoteConfigURL != "" {
		uRL, err := url.Parse(remoteConfigURL)
		if err == nil {
			e.cfg = config.NewRemoteConfig(uRL.Scheme, uRL.Host, uRL.Path, "toml")
			if e.cfg == nil {
				retry := 600
				for i := 0; i < retry; i++ {
					log.Infof( "wait load config")
					time.Sleep(time.Second)
					e.cfg = config.NewRemoteConfig(uRL.Scheme, uRL.Host, uRL.Path, "toml")
					if e.cfg != nil {
						break
					}
					if i+1 == retry {
						log.Fatalf( "cannot load config")
					}
				}
			}
			if localConfigURL != "" {
				dir := path.Dir(localConfigURL)
				_, err := os.Stat(dir)
				if os.IsNotExist(err) {
					_ = os.MkdirAll(dir, 0777)
				}
				_ = e.cfg.Save(localConfigURL)
			}
		}
	}
	if e.cfg == nil {
		if localConfigURL == "" {
			if _, err := os.Stat("conf-file.toml"); err == nil {
				localConfigURL = "conf-file.toml"
			}
		}
		if localConfigURL != "" {
			e.cfg = config.NewConfig(localConfigURL)
		}
	}

	log.Config(e.cfg.GetString("application.name"),
		log.Level(e.cfg.GetString("log.level")),
		e.cfg.GetBool("log.console"),
		e.cfg.GetString("log.filename"),
		func() log.LogFormat {
			if e.cfg.GetString("log.format") == "text" {
				return log.LogFormatText
			} else {
				return log.LogFormatJSON
			}
		}(),
		e.cfg.GetInt("log.maxSize"),
		e.cfg.GetInt("log.maxBackups"),
		e.cfg.GetInt("log.maxAge"),
	)

	log.Infof( "log level:%s", e.cfg.GetString("log.level"))
	return nil
}

func (e *Executor) Start() error {
	e.s.Main(e.cfg)
	return nil
}

func (e *Executor) Stop() error {
	if e.s != nil {
		e.s.Exit()
	}
	log.Flush()
	if e.cfg != nil {
		e.cfg.Exit()
	}
	return nil
}
