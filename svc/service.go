package svc

import (
	"github.com/brickman-source/golang-utilities/config"
)

type Service interface {
	GetName() string
	Main(config *config.Config)
	Exit()
}
