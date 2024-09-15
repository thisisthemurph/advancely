package errs

import "errors"

// IsOne reports whether any error in the `err` param's tree matches any of the targets.
func IsOne(err error, targets ...error) bool {
	for _, t := range targets {
		if errors.Is(err, t) {
			return true
		}
	}
	return false
}
