//go:build linux

package middleware

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// getCPULoad 读取 /proc/stat 采样两次，计算 CPU 使用率百分比（0~100）。
// 仅在 Linux 上编译。
func getCPULoad(interval time.Duration) float64 {
	idle0, total0 := readProcStat()
	time.Sleep(interval / 2)
	idle1, total1 := readProcStat()

	totalDiff := total1 - total0
	idleDiff := idle1 - idle0
	if totalDiff == 0 {
		return 0
	}
	usage := 100.0 * float64(totalDiff-idleDiff) / float64(totalDiff)
	if usage < 0 {
		return 0
	}
	if usage > 100 {
		return 100
	}
	return usage
}

// readProcStat 解析 /proc/stat 的第一行，返回 (idle, total) jiffies。
func readProcStat() (idle, total uint64) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0
	}
	line := strings.SplitN(string(data), "\n", 2)[0]
	fields := strings.Fields(line)
	if len(fields) < 5 || fields[0] != "cpu" {
		return 0, 0
	}
	for i, f := range fields[1:] {
		v, _ := strconv.ParseUint(f, 10, 64)
		total += v
		if i == 3 { // index 3 = idle
			idle = v
		}
	}
	return
}
