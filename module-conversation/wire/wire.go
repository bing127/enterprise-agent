package wire

import (
    "github.com/google/wire"
    "github.com/bing127/enterprise-agent/module-infra/middleware/es"
    "github.com/bing127/enterprise-agent/module-infra/middleware/weaviate"
    "github.com/bing127/enterprise-agent/module-infra/middleware/redis"
    "github.com/bing127/enterprise-agent/module-conversation/service"
)

var SuperSet = wire.NewSet(
    es.NewClient,
    weaviate.NewClient,
    redis.NewClient,
    service.NewConversationService,
)