package utils

import "errors"

// IPsChecker func
func IPsChecker(currentUserIP string, lastCachedUserIP string) error {
	if len(lastCachedUserIP) == 0 || currentUserIP == lastCachedUserIP {
		return nil
	}
	return errors.New("You may watch only on one device")
}
