package config

type CooldownConfig struct {
	OptCooldown            duration
	JoinLeaveCooldown      duration
	WarnViolationThreshold uint64
}
