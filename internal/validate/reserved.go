package validate

// ReservedKeyList returns a copy of all keys reserved for internal use.
// Useful for display in help text or shell completions.
func ReservedKeyList() []string {
	keys := make([]string, 0, len(reservedKeys))
	for k := range reservedKeys {
		keys = append(keys, k)
	}
	return keys
}
