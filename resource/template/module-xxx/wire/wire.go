package wire

import (
	"github.com/bing127/enterprise-agent/module-xxx/service"
	"github.com/google/wire"
)

// SuperSet 包含 module-xxx 所有可复用的 Wire Provider。
var SuperSet = wire.NewSet(
	service.NewXxxService,
)
