// Package encrypt provides AES-256-GCM encryption helpers for securing
// environment variable values at rest inside envoy configuration files.
//
// Encrypted values are stored with an "enc:" prefix followed by a
// base64-encoded nonce+ciphertext blob so they can be round-tripped
// through the standard config load/save cycle without modification.
//
// Usage:
//
//	enc, err := encrypt.EncryptValue(passphrase, "my-secret")
//	dec, err := encrypt.DecryptValue(passphrase, enc)
//	err       := encrypt.EncryptTarget(cfg, "production", passphrase)
//	err       := encrypt.DecryptTarget(cfg, "production", passphrase)
package encrypt
