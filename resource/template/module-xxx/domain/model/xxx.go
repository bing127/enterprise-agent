package model

import "time"

// Xxx 实体定义。
type Xxx struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Status    int8      `json:"status"` // 0=禁用 1=启用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
