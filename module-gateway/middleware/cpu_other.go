//go:build !linux

package middleware

import "time"

// getCPULoad 在非 Linux 系统（macOS / Windows）上返回 0，
// 自适应降载将仅基于在途请求数工作。
// 生产环境（Linux）会使用 cpu_linux.go 中读取 /proc/stat 的实现。
func getCPULoad(_ time.Duration) float64 {
	return 0
}
