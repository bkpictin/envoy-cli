// Package group lets users organise environment variable keys into named
// logical groups within a deployment target.
//
// Groups are purely organisational — they do not affect how values are
// exported or resolved. A key may belong to multiple groups simultaneously.
//
// Typical usage:
//
//	group.Create(cfg, "production", "database")
//	group.AddKey(cfg, "production", "database", "DB_HOST")
//	group.AddKey(cfg, "production", "database", "DB_PORT")
//	keys, _ := group.GetKeys(cfg, "production", "database")
package group
