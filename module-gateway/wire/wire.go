package wire

import (
	corewire "github.com/bing127/enterprise-agent/module-core/wire"
	"github.com/bing127/enterprise-agent/module-infra/middleware/es"
	"github.com/bing127/enterprise-agent/module-infra/middleware/redis"
	"github.com/bing127/enterprise-agent/module-infra/middleware/weaviate"
	"github.com/google/wire"
)

var SuperSet = wire.NewSet(
	es.NewClient,
	weaviate.NewClient,
	redis.NewClient,
	corewire.SuperSet,
)
