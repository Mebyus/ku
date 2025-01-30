package main

import (
	"fmt"
)

// Value represents an SSA value (e.g., a variable or constant).
type Value struct {
	ID   int
	Name string
}

// BasicBlock represents a basic block in the control flow graph.
type BasicBlock struct {
	ID      int
	Preds   []*BasicBlock // Predecessor blocks
	Succs   []*BasicBlock // Successor blocks
	Args    []*Value      // Arguments (values passed from predecessors)
	Insts   []Instruction // Instructions in the block
	OutArgs []*Value      // Outgoing arguments (values passed to successor blocks)
}

// Instruction represents an SSA instruction.
type Instruction interface {
	String() string
}

// BinaryOp represents a binary operation (e.g., add, sub).
type BinaryOp struct {
	Op    string
	Left  *Value
	Right *Value
	Dest  *Value
}

func (b *BinaryOp) String() string {
	return fmt.Sprintf("%s = %s %s, %s", b.Dest.Name, b.Op, b.Left.Name, b.Right.Name)
}

// Branch represents a conditional branch instruction.
type Branch struct {
	Cond      *Value      // Condition value
	TrueSucc  *BasicBlock // Successor if condition is true
	FalseSucc *BasicBlock // Successor if condition is false
}

func (b *Branch) String() string {
	return fmt.Sprintf("br %s, Block %d, Block %d", b.Cond.Name, b.TrueSucc.ID, b.FalseSucc.ID)
}

// Jump represents an unconditional jump instruction.
type Jump struct {
	Target *BasicBlock // Target block
}

func (j *Jump) String() string {
	return fmt.Sprintf("jump Block %d", j.Target.ID)
}

// NewBasicBlock creates a new basic block.
func NewBasicBlock(id int) *BasicBlock {
	return &BasicBlock{
		ID:      id,
		Preds:   []*BasicBlock{},
		Succs:   []*BasicBlock{},
		Args:    []*Value{},
		Insts:   []Instruction{},
		OutArgs: []*Value{},
	}
}

// AddBinaryOp adds a binary operation to the basic block.
func (b *BasicBlock) AddBinaryOp(op string, left, right, dest *Value) {
	binOp := &BinaryOp{Op: op, Left: left, Right: right, Dest: dest}
	b.Insts = append(b.Insts, binOp)
}

// AddBranch adds a conditional branch instruction to the basic block.
func (b *BasicBlock) AddBranch(cond *Value, trueSucc, falseSucc *BasicBlock) {
	branch := &Branch{Cond: cond, TrueSucc: trueSucc, FalseSucc: falseSucc}
	b.Insts = append(b.Insts, branch)
}

// AddJump adds an unconditional jump instruction to the basic block.
func (b *BasicBlock) AddJump(target *BasicBlock) {
	jump := &Jump{Target: target}
	b.Insts = append(b.Insts, jump)
}

// AddOutArg adds an outgoing argument to the basic block.
func (b *BasicBlock) AddOutArg(arg *Value) {
	b.OutArgs = append(b.OutArgs, arg)
}

// PrintBlock prints the contents of a basic block.
func PrintBlock(b *BasicBlock) {
	fmt.Printf("Block %d:\n", b.ID)
	fmt.Println("  Args:")
	for _, arg := range b.Args {
		fmt.Printf("    %s\n", arg.Name)
	}
	fmt.Println("  Outgoing Args:")
	for _, outArg := range b.OutArgs {
		fmt.Printf("    %s\n", outArg.Name)
	}
	fmt.Println("  Instructions:")
	for _, inst := range b.Insts {
		fmt.Printf("    %s\n", inst.String())
	}
}

func main() {
	// Create some basic blocks
	entry := NewBasicBlock(0)
	thenBlock := NewBasicBlock(1)
	elseBlock := NewBasicBlock(2)
	mergeBlock := NewBasicBlock(3)

	// Create some SSA values
	x := &Value{ID: 1, Name: "x"}
	y := &Value{ID: 2, Name: "y"}
	z := &Value{ID: 3, Name: "z"}
	cond := &Value{ID: 4, Name: "cond"}
	resultThen := &Value{ID: 5, Name: "result.then"}
	resultElse := &Value{ID: 6, Name: "result.else"}
	resultMerge := &Value{ID: 7, Name: "result.merge"}

	// Simulate control flow
	entry.Succs = []*BasicBlock{thenBlock, elseBlock}
	thenBlock.Preds = []*BasicBlock{entry}
	elseBlock.Preds = []*BasicBlock{entry}
	mergeBlock.Preds = []*BasicBlock{thenBlock, elseBlock}

	// Add instructions to the entry block
	entry.AddBinaryOp("add", x, y, z)
	entry.AddBinaryOp("eq", z, y, cond)         // cond = (z == y)
	entry.AddBranch(cond, thenBlock, elseBlock) // Branch based on cond

	// Pass z and y as arguments to thenBlock and elseBlock
	thenBlock.Args = []*Value{z, y}
	elseBlock.Args = []*Value{z, y}

	// In the thenBlock, compute `z + y` and assign it to `resultThen`
	thenBlock.AddBinaryOp("add", thenBlock.Args[0], thenBlock.Args[1], resultThen)
	thenBlock.AddOutArg(resultThen) // Pass `resultThen` to the merge block
	thenBlock.AddJump(mergeBlock)   // Jump to merge block

	// In the elseBlock, compute `z - y` and assign it to `resultElse`
	elseBlock.AddBinaryOp("sub", elseBlock.Args[0], elseBlock.Args[1], resultElse)
	elseBlock.AddOutArg(resultElse) // Pass `resultElse` to the merge block
	elseBlock.AddJump(mergeBlock)   // Jump to merge block

	// In the mergeBlock, use a single argument to represent the merged value
	mergeBlock.Args = []*Value{resultMerge}

	// The merged result is `resultMerge`, which is implicitly defined by the basic block arguments
	// (This is where the "phi node" logic would happen in a real compiler)
	fmt.Println("Implicit Phi Node in Merge Block:")
	fmt.Printf("  %s = φ(%s, %s)\n", resultMerge.Name, resultThen.Name, resultElse.Name)

	// Print the blocks
	PrintBlock(entry)
	PrintBlock(thenBlock)
	PrintBlock(elseBlock)
	PrintBlock(mergeBlock)
}
