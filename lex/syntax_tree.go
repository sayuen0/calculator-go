package lex

import "fmt"

type Value float64

// 評価できるものを式と呼ぶ
type Expr interface {
	Eval() Value
}

func (e Value) Eval() Value {
	return e
}

//  単項演算子
type Op1 struct {
	code rune
	expr Expr
}

func newOp1(code rune, e Expr) Expr {
	return &Op1{code, e}
}

// 単項演算子の評価
func (e *Op1) Eval() Value {
	v := e.expr.Eval()
	if e.code == '-' {
		v = -v
	}
	return v
}

// 二項演算子
type Op2 struct {
	code        rune
	left, right Expr
}

func newOp2(code rune, left, right Expr) Expr {
	return &Op2{code, left, right}
}

// 二項演算子の評価
func (e *Op2) Eval() Value {
	x := e.left.Eval()
	y := e.right.Eval()
	switch e.code {
	case '+':
		return x + y
	case '-':
		return x - y
	case '*':
		return x * y
	case '/':
		return x / y
	default:
		panic(fmt.Errorf("invalid op code"))
	}
}
