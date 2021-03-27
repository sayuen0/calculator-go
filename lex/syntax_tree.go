package lex

import "fmt"

type Value float64

// 局所変数の環境
type Env struct {
	name Variable
	val  Value
	next *Env
}

func newEnv(name Variable, val Value, next *Env) *Env {
	return &Env{name, val, next}
}

//  変数束縛
func makeBinding(xs []Variable, es []Expr, env *Env) *Env {
	var env1 *Env
	for i := 0; i < len(xs); i++ {
		env1 = newEnv(xs[i], es[i].Eval(env), env1)
	}
	return env1
}

// 構文木の型
type Expr interface {
	Eval(*Env) Value
}

// 局所変数の参照
func lookUp(name Variable, env *Env) (Value, bool) {
	for ; env != nil; env = env.next {
		if name == env.name {
			return env.val, true
		}
	}
	return 0.0, false
}

// 局所変数の更新
func update(name Variable, val Value, env *Env) bool {
	for ; env != nil; env = env.next {
		if name == env.name {
			env.val = val
			return true
		}
	}
	return false
}

func (e Value) Eval(env *Env) Value {
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
func (e *Op1) Eval(env *Env) Value {
	v := e.expr.Eval(env)
	switch e.code {
	case '-':
		return -v
	case '+':
		return v
	case NOT:
		return boolToValue(isFalse(v))
	default:
		panic(fmt.Errorf("invalid Op1 code"))
	}
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
func (e *Op2) Eval(env *Env) Value {
	x := e.left.Eval(env)
	y := e.right.Eval(env)
	switch e.code {
	case '+':
		return x + y
	case '-':
		return x - y
	case '*':
		return x * y
	case '/':
		return x / y
	case EQ:
		return boolToValue(x == y)
	case NE:
		return boolToValue(x != y)
	case LT:
		return boolToValue(x < y)
	case GT:
		return boolToValue(x > y)
	case LE:
		return boolToValue(x <= y)
	case GE:
		return boolToValue(x >= y)
	default:
		panic(fmt.Errorf("invalid op code"))
	}
}
