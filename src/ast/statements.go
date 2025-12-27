package ast

import "github.com/caelondev/monkey/src/token"

type LetStatement struct {
	Token token.Token // LET Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Identifier struct {
	Token token.Token // IDENTIFIER Token
	Value string
}

func (i *Identifier) statementNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
