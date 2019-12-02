package generic

import (
	"fmt"
	"strings"
)

// CheckHostname performs minimal validation to see if the hostname string could be a valid hostname. If the validation fails, CheckHostname returns an error with a message that explains the validation failure. If validation succeeds, CheckHostname returns nil.
func CheckHostname(hostname string) error {
	// Perform only minimal validation of the hostname, since this value should be validated downstream by actual DNS resolution (and, thorough checking would be very complicated... e.g., see https://stackoverflow.com/a/26618995).

	// Check length
	if length := len(hostname); length < 2 || length > 264 {
		return fmt.Errorf("invalid hostname length (hostname: '%s', length: %d)", hostname, length)
	}

	// Check to ensure two or more non-zero length domain components
	components := strings.Split(hostname, ".")
	if componentsCount := len(components); componentsCount < 2 {
		return fmt.Errorf("need 2 or more domain components (hostname: '%s', components: %d)", hostname, componentsCount)
	}

	for i, component := range components {
		if componentLength := len(component); componentLength == 0 {
			return fmt.Errorf("hostname cannot a zero-length domain component (hostname: '%s', component index with zero-length: %d)", hostname, i)
		}
	}

	return nil
}
