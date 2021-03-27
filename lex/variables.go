package lex

import (
	"fmt"
	"math"
)

// 変数
type Variable string

// 大域的な環境
var globalEnv = make(map[Variable]Value)

// テスト用
func resetGlobal() {
	globalEnv = make(map[Variable]Value)
}

// 変数の評価
func (v Variable) Eval() Value {
	val, ok := globalEnv[v]
	if !ok {
		panic(fmt.Errorf("unbound variable: %v", v))
	}
	return val
}

// 代入演算子
type Agn struct {
	name Variable
	expr Expr
}

func newAgn(v Variable, e Expr) *Agn {
	return &Agn{v, e}
}

// 代入演算子の評価
func (a *Agn) Eval() Value {
	val := a.expr.Eval()
	globalEnv[a.name] = val
	return val
}

// 組み込み関数 引数の個数を持てるもの
type Func interface {
	Argc() int // 引数の個数
}

type Func1 func(float64) float64

func (f Func1) Argc() int {
	return 1
}

type Func2 func(float64, float64) float64

func (f Func2) Argc() int {
	return 2
}

// 組み込み関数の構文木
type App struct {
	fn Func
	xs []Expr
}

func newApp(fn Func, xs []Expr) *App {
	return &App{fn: fn, xs: xs}
}

// 組み込み関数の評価
func (a *App) Eval() Value {
	switch f := a.fn.(type) {
	case Func1:
		fn := f
		x := float64(a.xs[0].Eval())
		return Value(fn(x))
	case Func2:
		fn := f
		x := float64(a.xs[0].Eval())
		y := float64(a.xs[1].Eval())
		return Value(fn(x, y))
	default:
		panic(fmt.Errorf("function Eval error"))
	}
}

// 組み込み関数表
var funcTable = make(map[string]Func)

// 組み込み関数初期化
func InitFunc() {
	funcTable["sqrt"] = Func1(math.Sqrt)
	funcTable["sin"] = Func1(math.Sin)
	funcTable["cos"] = Func1(math.Cos)
	funcTable["tan"] = Func1(math.Tan)
	funcTable["sinh"] = Func1(math.Sinh)
	funcTable["cosh"] = Func1(math.Cosh)
	funcTable["tanh"] = Func1(math.Tanh)
	funcTable["asin"] = Func1(math.Asin)
	funcTable["acos"] = Func1(math.Acos)
	funcTable["atan"] = Func1(math.Atan)
	funcTable["atan2"] = Func2(math.Atan2)
	funcTable["exp"] = Func1(math.Exp)
	funcTable["pow"] = Func2(math.Pow)
	funcTable["log"] = Func1(math.Log)
	funcTable["log10"] = Func1(math.Log10)
	funcTable["log2"] = Func1(math.Log2)
	funcTable["abs"] = Func1(math.Abs)
}
