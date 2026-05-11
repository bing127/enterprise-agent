# Enterprise Agent

This project is designed to build an enterprise-grade intelligent agent using Eino and Hertz. The architecture is structured to ensure scalability and maintainability as the project grows.

## Project Structure

- **go.work**: Defines the workspace for managing multiple Go modules.
- **go.work.sum**: Contains checksums for all dependencies in the workspace.

### Infrastructure

- **infra/go.mod**: Dependency management for the infrastructure module.
- **infra/middleware**: Contains middleware implementations.
  - **es/client.go**: Client interactions with Elasticsearch.
  - **weaviate/client.go**: Client interactions with Weaviate.
  - **redis/client.go**: Client interactions with Redis.
- **infra/wire/wire.go**: Dependency injection configuration using Wire.

### Agent Module

- **agent/go.mod**: Dependency management for the agent module.
- **agent/wire**: Dependency injection configuration for the agent module.
  - **wire.go**: Wire configuration.
  - **wire_gen.go**: Generated code for dependency injection.
- **agent/domain**: Contains domain models and repositories.
  - **model/agent.go**: Defines the agent model structure.
  - **repo/agent_repo.go**: Repository interface for agent-related data.
- **agent/service/agent_service.go**: Implements business logic for the agent.
- **agent/handler/agent_handler.go**: Handles HTTP requests related to the agent.

### Knowledge Module

- **knowledge/go.mod**: Dependency management for the knowledge module.
- **knowledge/wire**: Dependency injection configuration for the knowledge module.
  - **wire.go**: Wire configuration.
  - **wire_gen.go**: Generated code for dependency injection.
- **knowledge/domain**: Contains domain models and repositories.
  - **model/knowledge.go**: Defines the knowledge model structure.
  - **repo/knowledge_repo.go**: Repository interface for knowledge-related data.
- **knowledge/service/knowledge_service.go**: Implements business logic for knowledge.
- **knowledge/handler/knowledge_handler.go**: Handles HTTP requests related to knowledge.

### Conversation Module

- **conversation/go.mod**: Dependency management for the conversation module.
- **conversation/wire**: Dependency injection configuration for the conversation module.
  - **wire.go**: Wire configuration.
  - **wire_gen.go**: Generated code for dependency injection.
- **conversation/domain**: Contains domain models and repositories.
  - **model/conversation.go**: Defines the conversation model structure.
  - **repo/conversation_repo.go**: Repository interface for conversation-related data.
- **conversation/service/conversation_service.go**: Implements business logic for conversation.
- **conversation/handler/conversation_handler.go**: Handles HTTP requests related to conversation.

### Gateway Module

- **gateway/go.mod**: Dependency management for the gateway module.
- **gateway/wire**: Dependency injection configuration for the gateway module.
  - **wire.go**: Wire configuration.
  - **wire_gen.go**: Generated code for dependency injection.
- **gateway/middleware**: Contains middleware implementations.
  - **auth.go**: Authentication middleware.
  - **rate_limit.go**: Rate limiting middleware.
  - **trace.go**: Request tracing middleware.
- **gateway/router/router.go**: Defines routing configuration for the gateway.

### Configuration Module

- **config/go.mod**: Dependency management for the configuration module.
- **config/config.go**: Implements loading and managing configurations.

## Getting Started

To get started with the project, ensure you have Go installed and set up your environment. You can then run the following commands to initialize the workspace and install dependencies:

```bash
go work init
go mod tidy
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any suggestions or improvements.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.