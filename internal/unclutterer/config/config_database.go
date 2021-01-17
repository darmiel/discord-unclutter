package config

type DatabaseConfig struct {
	AllowGhostPingBlocking bool
	GhostPingBlockDefault  bool
	DatabaseBucketName     string
	DatabasePath           string
}
