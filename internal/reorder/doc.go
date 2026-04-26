// Package reorder provides utilities for reordering environment variable keys
// within a target. Keys can be sorted alphabetically or arranged in a custom
// order specified by the caller.
//
// Alphabetical reordering sorts all keys in a target using standard lexicographic
// ordering, which improves readability and makes diffs easier to review.
//
// Custom reordering places a specified subset of keys at the top of the list in
// the given order, with any remaining keys appended afterward in their original
// relative order.
package reorder
