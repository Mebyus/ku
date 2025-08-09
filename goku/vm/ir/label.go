package ir

// Label contains label name in integer form.
//
// Directly corresponds to label entry index inside list of
// all program text labels.
type Label uint32

// Place represents label placement operation in text segment.
type Place struct {
	nodeAtom

	Label Label
}
