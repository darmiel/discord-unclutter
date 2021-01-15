package main

import "time"

const cooldown = 5 * time.Second

var cooldownUsers = make(map[string]time.Time)
var cooldownViolations = make(map[string]uint64)

func cooldownUser(userID string) {
	cooldownUsers[userID] = time.Now().Add(cooldown)
	cooldownViolations[userID] = 0
}

func checkAndUpdateCooldown(userID string) (cooldowned bool, violations uint64) {
	dura, ok := cooldownUsers[userID]
	if !ok {
		cooldownUser(userID)
		return false, 0
	}

	cooldowned = time.Now().Before(dura)
	if !cooldowned {
		cooldownUser(userID)
	} else {
		violations, ok = cooldownViolations[userID]
		if !ok {
			violations = 0
		}
		violations++
		cooldownViolations[userID] = violations
	}

	return
}
