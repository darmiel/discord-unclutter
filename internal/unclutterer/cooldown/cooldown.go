package cooldown

import "time"

const cooldown = 5 * time.Second

var (
	cooldowns  = make(map[string]time.Time)
	violations = make(map[string]uint64)
)

func SetOnCooldown(obj string) {
	cooldowns[obj] = time.Now().Add(cooldown)
	violations[obj] = 0
}

func IsOnCooldown(obj string) (cooldowned bool, vl uint64) {
	dura, ok := cooldowns[obj]

	// if not cooldowned
	// -> cooldown
	if !ok {
		SetOnCooldown(obj)
		return false, 0
	}

	cooldowned = time.Now().Before(dura)
	if !cooldowned {
		SetOnCooldown(obj)
	} else {
		vl, ok = violations[obj]
		if !ok {
			vl = 0
		}
		vl++

		// update violations
		violations[obj] = vl
	}

	return
}
