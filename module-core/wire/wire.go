package wire

import (
	"github.com/bing127/enterprise-agent/module-core/handler"
	"github.com/bing127/enterprise-agent/module-core/service"
	"github.com/google/wire"
)

// SuperSet 包含 module-core 所有可复用的 Wire Provider。
var SuperSet = wire.NewSet(
	service.NewAuthService,
	handler.NewAuthHandler,
)
