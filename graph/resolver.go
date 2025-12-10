package graph

import (
	"github.com/rohit21755/gg_server.git/ws"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	DB    *gorm.DB
	WsHub *ws.Hub
}
