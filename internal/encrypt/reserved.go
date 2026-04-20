package encrypt

// SensitiveKeyHints is a list of key substrings that are commonly associated
// with sensitive values. These can be used by tooling to suggest encryption
// to the user but are not enforced automatically.
var SensitiveKeyHints = []string{
	"PASSWORD",
	"PASS",
	"SECRET",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
}

// IsSensitiveKey returns true when the key name contains any of the
// SensitiveKeyHints substrings (case-insensitive comparison via upper-case).
func IsSensitiveKey(key string) bool {
	upper := toUpper(key)
	for _, hint := range SensitiveKeyHints {
		if contains(upper, hint) {
			return true
		}
	}
	return false
}

func toUpper(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 32
		}
		b[i] = c
	}
	return string(b)
}

func contains(s, sub string) bool {
	if len(sub) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
