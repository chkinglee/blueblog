// Package trace
// @Author      : lilinzhen
// @Time        : 2022/4/11 17:26:15
// @Description :
package trace

type RabbitMQ struct {
	Timestamp   string      `json:"timestamp"`      // 时间，格式：2006-01-02 15:04:05
	Exchange    string      `json:"exchange"`       // Exchange
	RoutingKey  string      `json:"routing_key"`    // RoutingKey
	Queue       string      `json:"queue"`          // Queue
	Data        interface{} `json:"data,omitempty"` // Data
	CostSeconds float64     `json:"cost_seconds"`   // 执行时间(单位秒)
}
