package kub

// Target describes build target.
type Target struct {
	Arch Arch
	OS   OS
}

// Build describes build conditions.
type Build struct {
	Target Target
	Type   Type
	Goal   Goal
}
