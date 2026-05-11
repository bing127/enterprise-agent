// Package middleware — 自适应降载中间件。
//
// 当 CPU 使用率超过 maxCPU% 且在途请求数超过 maxInFlight 时，
// 以 shed 概率拒绝新请求，返回 503，避免雪崩。
package middleware

import (
	"context"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/bing127/enterprise-agent/module-pkg/logger"
	"github.com/bing127/enterprise-agent/module-pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
	"go.uber.org/zap"
)

// SheddingOptions 自适应降载配置。
type SheddingOptions struct {
	// MaxCPU CPU 使用率阈值（0~100），超过时触发降载评估。
	// 非 Linux 系统该字段无效（始终视为 0%）。
	MaxCPU float64

	// MaxInFlight 在途请求数上限，超过时触发降载评估。
	MaxInFlight int64

	// SampleInterval CPU 采样间隔（默认 500ms）。
	SampleInterval time.Duration
}

type sheddingState struct {
	inFlight  atomic.Int64
	cpuLoad   atomic.Value // float64
	shedCount atomic.Int64
}

// AdaptiveSheddingMiddleware 自适应降载中间件。
//
// 同时满足以下两个条件时开始以概率方式拒绝请求：
//  1. CPU 使用率 > opts.MaxCPU
//  2. 在途请求数 > opts.MaxInFlight
//
// 拒绝概率随过载程度线性增长，最高 90%（保留 10% 用于探测恢复）。
func AdaptiveSheddingMiddleware(opts SheddingOptions) app.HandlerFunc {
	if opts.SampleInterval <= 0 {
		opts.SampleInterval = 500 * time.Millisecond
	}

	st := &sheddingState{}
	st.cpuLoad.Store(float64(0))

	log := logger.L().Named("shedding")

	// 后台 goroutine 定期采样 CPU
	go func() {
		for {
			time.Sleep(opts.SampleInterval)
			load := getCPULoad(opts.SampleInterval)
			st.cpuLoad.Store(load)
		}
	}()

	return func(ctx context.Context, c *app.RequestContext) {
		inFlight := st.inFlight.Add(1)
		defer st.inFlight.Add(-1)

		cpuLoad := st.cpuLoad.Load().(float64)

		// 双阈值判断：CPU 过高 AND 在途请求过多
		if cpuLoad > opts.MaxCPU && inFlight > opts.MaxInFlight {
			// 过载比例（取 CPU 和在途请求中较严重的一方）
			cpuOverload := (cpuLoad - opts.MaxCPU) / (100 - opts.MaxCPU)
			inFlightOverload := float64(inFlight-opts.MaxInFlight) / float64(opts.MaxInFlight)
			overload := max64(cpuOverload, inFlightOverload)

			// 拒绝概率：[0, 0.9]
			shedProb := overload * 0.9
			if rand.Float64() < shedProb {
				shed := st.shedCount.Add(1)
				log.Warn("request shed",
					zap.Float64("cpu_load", cpuLoad),
					zap.Int64("in_flight", inFlight),
					zap.Float64("shed_prob", shedProb),
					zap.Int64("total_shed", shed),
				)
				response.ServiceUnavailable(ctx, c, "server is overloaded, please retry later")
				c.Abort()
				return
			}
		}

		c.Next(ctx)
	}
}

func max64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
