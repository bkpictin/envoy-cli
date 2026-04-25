// Package pin implements key-pinning for envoy-cli targets.
//
// A pinned key is protected from being silently overwritten by bulk
// operations such as sync, merge, copy, and promote. Any operation that
// would modify a pinned key must be invoked with the --force flag, or it
// will skip that key and emit a warning.
//
// Pinned state is stored inline within the target's var map using a
// reserved "__pinned__<KEY>" sentinel entry so that no additional config
// section is required.
package pin
