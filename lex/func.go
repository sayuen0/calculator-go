package lex

import (
	"fmt"
	"text/scanner"
)

// ユーザ定義関数
type FuncU struct {
	name string
	xs   []Variable
	body Expr
}

func newFuncU(name string, xs []Variable, body Expr) *FuncU {
	return &FuncU{name, xs, body}
}

func (f *FuncU) Argc() int {
	return len(f.xs)
}

// ユーザ関数の定義
func defineFunc(lex *Lex) {
	lex.getToken()
	if lex.Token != scanner.Ident {
		panic(fmt.Errorf("invalid define form"))
	}
	name := lex.TokenText()
	lex.getToken()
	xs := getParameter(lex)
	body := expression(lex)
	if lex.Token != END {
		panic(fmt.Errorf("'end' expected"))
	}
	v, ok := funcTable[name]
	if ok {
		switch f := v.(type) {
		case *FuncU:
			if len(f.xs) != len(xs) {
				panic(fmt.Errorf("wrong number of arguments: %v", name))
			}
			f.xs = xs
			f.body = body
		default:
			panic(fmt.Errorf("%v is built-in function", name))
		}
	} else {
		funcTable[name] = newFuncU(name, xs, body)
	}
	fmt.Println(name)
}

// 仮引数の取得
func getParameter(lex *Lex) []Variable {
	e := make([]Variable, 0)
	if lex.Token != '(' {
		panic(fmt.Errorf("'(' expected"))
	}
	lex.getToken()
	if lex.Token == ')' {
		lex.getToken()
		return e
	}
	for {
		if lex.Token == scanner.Ident {
			e = append(e, Variable(lex.TokenText()))
			lex.getToken()
			switch lex.Token {
			case ')':
				lex.getToken()
				return e
			case ',':
				lex.getToken()
			default:
				panic(fmt.Errorf("unexpected token in parameter list"))
			}
		} else {
			panic(fmt.Errorf("unexpected token in parameter list"))
		}
	}
}
