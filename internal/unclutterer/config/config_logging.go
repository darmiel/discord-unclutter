package config

type LoggingConfig struct {
	VerbosityLevel uint

	LogChat                        bool
	LogUserJoin                    bool
	LogUserLeave                   bool
	LogCooldown                    bool
	LogExcessiveCooldownViolations bool

	LogOptOutError   bool
	LogOptInError    bool
	LogOptOutSuccess bool
	LogOptInSuccess  bool

	PingBindAddress string
}
