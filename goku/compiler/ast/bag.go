package ast

// RegBag represents top node which registers a type for bag implementation.
type RegBag struct {
	// Name of bag implementation being registered.
	Name Word

	// Bag name under which type is registered.
	BagName Word

	// Always not empty.
	Type TypeSpec

	Tab []BagTabField
}

type BagTabField struct {
	Name Word
	Fun  Word
}
