package config

import (
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/pkg/websocket"
	"go.uber.org/zap"
	"log"
)

// for ease of development we will switch to a closure pattern of dependency injection
// The reason I made this decision is because of two reasons
//		- Using method based dependency injection will not work because my handlers will be split in many packages.
//		I do not want to push in all of my handlers in one file only
//
//		- Using an external library introduces 'magic' which I do not want.

type Config struct {
	BaseURL  string
	HttpPort int
	Db       struct {
		Dsn         string
		Automigrate bool
	}
	Jwt struct {
		SecretKey string
	}
	Version bool
}

type Application struct {
	Config     Config
	Client     *ent.Client
	Logger     *zap.SugaredLogger
	MainLogger *zap.Logger
	StdLogger  *log.Logger
	Manager    *websocket.Manager
}
