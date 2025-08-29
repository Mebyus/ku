package vm

type ErrorCode uint32

const (
	// Execution reached end of program text segment.
	ErrorTextEnd ErrorCode = 1 + iota

	// Execution reached trap instruction.
	ErrorTrap

	ErrorBadOpcode

	ErrorBadSegment

	ErrorBadSpecialRegister

	ErrorBadRegister

	ErrorBadJumpAddress

	ErrorBadCallAddress

	ErrorBadVariant

	ErrorBadJumpFlag

	ErrorEmptyFrameStack

	ErrorReadOnlyRegister

	ErrorBadInstructionDataLength

	ErrorNonTextJump
)

type RuntimeError struct {
	Aux  uint64
	Code ErrorCode
}

func (r *RuntimeError) Error() string {
	return ""
}
