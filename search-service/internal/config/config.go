package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

// GetConf :nodoc:
func GetConf() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
	viper.AddConfigPath("./../../..")
	viper.SetConfigName("config")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Warningf("%v", err)
	}
}

// Env :nodoc:
func Env() string {
	return viper.GetString("env")
}

// LogLevel :nodoc:
func LogLevel() string {
	return viper.GetString("log_level")
}

// HTTPPort :nodoc:
func HTTPPort() string {
	return viper.GetString("search_service.ports.http")
}

// ElasticsearchHost :nodoc:
func ElasticsearchHost() string {
	return viper.GetString("elasticsearch.host")
}

// ElasticsearchMaxIdleConnections :nodoc:
func ElasticsearchMaxIdleConnections() int {
	if !viper.IsSet("elasticsearch.max_idle_connections") {
		return DefaultMaxElasticsearchIdleConnections
	}
	return viper.GetInt("elasticsearch.max_idle_connections")
}

// ElasticsearchMaxConnsPerHost :nodoc:
func ElasticsearchMaxConnsPerHost() int {
	if !viper.IsSet("elasticsearch.max_conns_per_host") {
		return DefaultMaxElasticsearchConnsPerHost
	}
	return viper.GetInt("elasticsearch.max_conns_per_host")
}

// ElasticsearchSetSniff :nodoc:
func ElasticsearchSetSniff() bool {
	if !viper.IsSet("elasticsearch.set_sniff") {
		return DefaultElasticsearchSetSniff
	}

	return viper.GetBool("elasticsearch.set_sniff")
}

// ElasticsearchSetHealthcheck :nodoc:
func ElasticsearchSetHealthcheck() bool {
	if !viper.IsSet("elasticsearch.set_health_check") {
		return DefaultElasticsearchSetHealthcheck
	}
	return viper.GetBool("elasticsearch.set_health_check")
}

// ElasticsearchHTTPTimeout :nodoc:
func ElasticsearchHTTPTimeout() time.Duration {
	return parseDuration(viper.GetString("elasticsearch.timeout"), DefaultElasticsearchHTTPTimeout)
}

// RedisDialTimeout :nodoc:
func RedisDialTimeout() time.Duration {
	cfg := viper.GetString("redis.dial_timeout")
	return parseDuration(cfg, 5*time.Second)
}

// RedisWriteTimeout :nodoc:
func RedisWriteTimeout() time.Duration {
	cfg := viper.GetString("redis.write_timeout")
	return parseDuration(cfg, 2*time.Second)
}

// RedisReadTimeout :nodoc:
func RedisReadTimeout() time.Duration {
	cfg := viper.GetString("redis.read_timeout")
	return parseDuration(cfg, 2*time.Second)
}

// RedisMaxIdleConn :nodoc:
func RedisMaxIdleConn() int {
	if viper.GetInt("redis.max_idle_conn") > 0 {
		return viper.GetInt("redis.max_idle_conn")
	}
	return 20
}

// RedisMaxActiveConn :nodoc:
func RedisMaxActiveConn() int {
	if viper.GetInt("redis.max_active_conn") > 0 {
		return viper.GetInt("redis.max_active_conn")
	}
	return 50
}

// ArticleServiceTarget :nodoc:
func ArticleServiceTarget() string {
	return viper.GetString("services.article_target")
}

//RPCServerTimeout :nodoc:
func RPCServerTimeout() time.Duration {
	if !viper.IsSet("rpc_server_timeout") {
		return time.Duration(DefaultRPCServerTimeout) * time.Millisecond
	}

	return time.Duration(viper.GetInt("rpc_server_timeout")) * time.Millisecond
}

//RPCClientTimeout :nodoc:
func RPCClientTimeout() time.Duration {
	if !viper.IsSet("rpc_client_timeout") {
		return time.Duration(DefaultRPCClientTimeout) * time.Millisecond
	}

	return time.Duration(viper.GetInt("rpc_client_timeout")) * time.Millisecond
}

// ServiceMaxConnPool :nodoc:
func ServiceMaxConnPool() int {
	if viper.GetInt("services.max_conn_pool") > 0 {
		return viper.GetInt("services.max_conn_pool")
	}
	return 200
}

// ServiceIdleConnPool :nodoc:
func ServiceIdleConnPool() int {
	if viper.GetInt("services.idle_conn_pool") > 0 {
		return viper.GetInt("services.idle_conn_pool")
	}
	return 100
}

// RedisPubSubHost :nodoc:
func RedisPubSubHost() string {
	return viper.GetString("redis.pub_sub")
}

func parseDuration(in string, defaultDuration time.Duration) time.Duration {
	dur, err := time.ParseDuration(in)
	if err != nil {
		return defaultDuration
	}
	return dur
}
