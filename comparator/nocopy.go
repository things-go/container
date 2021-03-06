package comparator

// noCopy may be embedded into structs which must not be copied
// after the first use.
//
// See https://golang.org/issues/8005#issuecomment-190753527
// for details.
type noCopy struct{} // nolint: unused

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {} // nolint: unused
func (*noCopy) Unlock() {} // nolint: unused
