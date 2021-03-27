package lex

import "fmt"

type Bgn struct {
	body []Expr
}

// begin
func newBgn(body []Expr) *Bgn {
	return &Bgn{body: body}
}

// begin式の処理
func getBody(lex *Lex) []Expr {
	body := make([]Expr, 0)
	for {
		body = append(body, expression(lex))
		switch lex.Token {
		case ',':
			lex.getToken()
		default:
			return body
		}
	}
}

func makeBegin(lex *Lex) Expr {
	if lex.Token == END {
		panic(fmt.Errorf("invalid begin form"))
	}
	body := getBody(lex)
	if lex.Token != END {
		panic(fmt.Errorf("'end' expected"))
	}
	lex.getToken()
	return newBgn(body)
}

// beginの評価
func (e *Bgn) Eval(env *Env) Value{
	var r Value
	for _, expr := range e.body{
		r = expr.Eval(env)
	}
	return r
}

// while
type Whl struct {
	testForm, body Expr
}

func newWhl(testForm, body Expr) *Whl {
	return &Whl{testForm, body}
}

// while式の処理
func makeWhile(lex *Lex) Expr {
	testForm := expression(lex)
	if lex.Token == DO {
		lex.getToken()
		return newWhl(testForm, makeBegin(lex))
	} else {
		panic(fmt.Errorf("'do' expected"))
	}
}

//whileの評価
func (e *Whl) Eval(env *Env) Value{
	for isTrue(e.testForm.Eval(env)){
		e.body.Eval(env)
	}
	return Value(0.0)
}

// let
type Let struct {
	vars []Variable
	vals []Expr
	body Expr
}

func newLet(vars []Variable, vals []Expr, body Expr) *Let {
	return &Let{vars: vars, vals: vals, body: body}
}

func makeLet(lex *Lex) Expr {
	vars := make([]Variable, 0)
	vals := make([]Expr, 0)
	for {

		e := expression(lex)
		a, ok := e.(*Agn)
		if !ok {
			panic(fmt.Errorf("let : invalid assign form"))
		}
		vars  =append(vars, a.name)
		vals = append(vals, a.expr)
		if lex.Token == IN{
			break
		} else if lex.Token != ','{
			panic(fmt.Errorf("let: ',' expected"))
		}
		lex.getToken()
	}
	lex.getToken()
	return newLet(vars, vals, makeBegin(lex))
}

// letの評価
func (e *Let) Eval(env *Env) Value{
	return e.body.Eval(addBinding(e.vars, e.vals, env))
}

// 局所変数を環境に追加
func addBinding(xs []Variable, es []Expr, env *Env) *Env{
	for i := 0 ; i < len(xs); i++{
		env = newEnv(xs[i], es[i].Eval(env), env)
	}
	return env
}