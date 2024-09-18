package abstract

// The frame is an abstract representation of a context/scope
// for instance, an entire source file, a function, or an if block
type Frame struct {
	Frames []Frame
	Fields []any
	Statements []Statement
}