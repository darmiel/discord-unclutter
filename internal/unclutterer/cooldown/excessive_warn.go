package cooldown

import (
	duconfig "github.com/darmiel/discord-unclutterer/internal/unclutterer/config"
	"log"
)

func CheckAndWarn(obj string, vl uint64, config *duconfig.Config) {
	if vl >= config.WarnViolationThreshold && config.LogExcessiveCooldownViolations {
		log.Println("  └ ⚠️ Obj", obj, "has a high amount of violations:", vl)
	}
}
