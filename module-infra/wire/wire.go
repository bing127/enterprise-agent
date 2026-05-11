package wire

import (
	"github.com/google/wire"
	"github.com/bing127/enterprise-agent/module-infra/middleware/es"
	"github.com/bing127/enterprise-agent/module-infra/middleware/weaviate"
	"github.com/bing127/enterprise-agent/module-infra/middleware/redis"
)

var InfraSet = wire.NewSet(
	es.NewClient,
	weaviate.NewClient,
	redis.NewClient,
)