module github.com/bing127/enterprise-agent/module-xxx

go 1.21

require (
	github.com/bing127/enterprise-agent/module-pkg v0.0.0
	github.com/cloudwego/hertz v0.10.4
	github.com/google/wire v0.7.0
	go.uber.org/zap v1.27.0
)

replace github.com/bing127/enterprise-agent/module-pkg => ../module-pkg
