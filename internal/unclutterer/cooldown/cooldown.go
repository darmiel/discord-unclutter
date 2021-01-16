package cooldown

import "time"

var (
	cooldowns  = make(map[string]time.Time)
	violations = make(map[string]uint64)
)

func SetOnCooldown(obj string, d time.Duration) {
	cooldowns[obj] = time.Now().Add(d)
	violations[obj] = 0
}

func IsOnCooldown(obj string, d time.Duration) (cooldowned bool, vl uint64) {
	dura, ok := cooldowns[obj]

	// if not cooldowned
	// -> cooldown
	if !ok {
		SetOnCooldown(obj, d)
		return false, 0
	}

	cooldowned = time.Now().Before(dura)
	if !cooldowned {
		SetOnCooldown(obj, d)
	} else {
		vl, ok = violations[obj]
		if !ok {
			vl = 0
		}
		vl++

		CheckAndWarn(obj, vl)

		// update violations
		violations[obj] = vl
	}

	return
}
