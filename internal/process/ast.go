package process

import (
	"cbeimers113/ilo/internal/constant"
	"fmt"
)

// All AST node types implement this interface for pretty printing
type ASTNode interface {
	String() string
}

type BinaryExpression struct {
	Left  ASTNode
	Op    Token
	Right ASTNode
}

type Variable struct {
	Name Token
}

type Literal struct {
	Value Token
}

type Assignment struct {
	Variable *Variable
	Value    ASTNode
}

type FuncDef struct {
	Name   Token
	Params []*Variable
	Body   []ASTNode
}

// String methods

func (b *BinaryExpression) String() string {
	return fmt.Sprintf("%sop(%s%s %s %s%s)%s", constant.ColRed, constant.ColReset, b.Left, b.Op.Literal, b.Right, constant.ColRed, constant.ColReset)
}

func (v *Variable) String() string {
	return fmt.Sprintf("%svar(%s%s%s)%s", constant.ColYellow, constant.ColReset, v.Name.Literal, constant.ColYellow, constant.ColReset)
}

func (l *Literal) String() string {
	return fmt.Sprintf("%slit(%s%s%s)%s", constant.ColYellow, constant.ColReset, l.Value.Literal, constant.ColYellow, constant.ColReset)
}

func (a *Assignment) String() string {
	return fmt.Sprintf("%sasgn(%s%s = %s%s)%s", constant.ColBlue, constant.ColReset, a.Variable, a.Value, constant.ColBlue, constant.ColReset)
}

func (f *FuncDef) String() string {
	body := ""
	for _, node := range f.Body {
		body += fmt.Sprintf("\t%s\n", node)
	}

	return fmt.Sprintf("%sfunc(%s%s(%v) {\n%s}%s)%s", constant.ColGreen, constant.ColReset, f.Name.Literal, f.Params, body, constant.ColGreen, constant.ColReset)
}
