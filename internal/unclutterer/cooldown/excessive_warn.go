package cooldown

import "log"

const Threshold = 4

func CheckAndWarn(obj string, vl uint64) {
	if vl >= Threshold {
		// TODO: Use from config
		log.Println("  └ ⚠️ Obj", obj, "has a high amount of violations:", vl)
	}
}
