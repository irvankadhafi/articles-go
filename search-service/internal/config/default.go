package config

import "time"

const (
	DefaultMaxElasticsearchIdleConnections = 2
	DefaultMaxElasticsearchConnsPerHost    = 10
	DefaultElasticsearchSetSniff           = false
	DefaultElasticsearchSetHealthcheck     = false
	DefaultElasticsearchHTTPTimeout        = 3 * time.Second

	DefaultRPCServerTimeout = 1000
	DefaultRPCClientTimeout = 1100
)
