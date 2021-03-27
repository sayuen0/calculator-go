package lex

import (
	"fmt"
	"os"
	"text/scanner"
)

func init() {
	initKeyTable()
}

// Lexとは、レキシカルアナライザ(字句解析プログラム)ジェネレータのこと
type Lex struct {
	scanner.Scanner
	Token rune
}

const (
	DEF = -(iota + 10)
	END
)

var keyTable = make(map[string]rune)

func initKeyTable() {
	keyTable["def"] = DEF
	keyTable["end"] = END
}

// 標準入力を1つ読み込んでruneを持つ
func (s *Lex) getToken() {
	s.Token = s.Scan()
	if s.Token == scanner.Ident {
		key, ok := keyTable[s.TokenText()]
		if ok {
			s.Token = key
		}
	}
}

// 引数の取得
func getArgs(lex *Lex) []Expr {
	e := make([]Expr, 0)
	if lex.Token != '(' {
		panic(fmt.Errorf("'(' expected"))
	}
	lex.getToken()
	if lex.Token == ')' {
		lex.getToken()
		return e
	}
	for {
		e = append(e, expression(lex))
		switch lex.Token {
		case ')':
			lex.getToken()
			return e
		case ',':
			lex.getToken()
		default:
			panic(fmt.Errorf("unexpected token in argument list"))
		}
	}
}

// factor: 因子
// 因子 = 数値 | ("+" | "-"), 因子 | "(" 式 ")".
func factor(lex *Lex) Expr {
	switch lex.Token {
	case '(':
		lex.getToken()
		e := expression(lex)
		if lex.Token != ')' {
			panic(fmt.Errorf("')' expected"))
		}
		lex.getToken()
		return e
	case '+':
		lex.getToken()
		return newOp1('+', factor(lex))
	case '-':
		lex.getToken()
		return newOp1('-', factor(lex))
	case scanner.Int, scanner.Float:
		var n float64
		fmt.Sscan(lex.TokenText(), &n)
		lex.getToken()
		return Value(n)
	case scanner.Ident:
		name := lex.TokenText()
		lex.getToken()
		if name == "quit" {
			panic(name)
		}
		v, ok := funcTable[name]
		if ok {
			xs := getArgs(lex)
			if len(xs) != v.Argc() {
				panic(fmt.Errorf("wrong number of argumnts: %v", name))
			}
			return newApp(v, xs)
		} else {
			return Variable(name)
		}
	default:
		panic(fmt.Errorf("unexpected token: %v", lex.TokenText()))
	}
}

// term: 項
// 項  = 因子 { ("*" | "/"), 因子 }.
func term(lex *Lex) Expr {
	e := factor(lex)
	for {
		switch lex.Token {
		case '*':
			lex.getToken()
			e = newOp2('*', e, factor(lex))
		case '/':
			lex.getToken()
			e = newOp2('/', e, factor(lex))
		default:
			return e
		}
	}
}

// 式
func expr1(lex *Lex) Expr {
	e := term(lex)
	for {
		switch lex.Token {
		case '+':
			lex.getToken()
			e = newOp2('+', e, term(lex))
		case '-':
			lex.getToken()
			e = newOp2('-', e, term(lex))
		default:
			return e
		}
	}
}

// expression: 式
// 式  = 項 { ("+" | "-"), 項 }.
func expression(lex *Lex) Expr {
	e := expr1(lex)
	if lex.Token == '=' {
		v, ok := e.(Variable)
		if ok {
			lex.getToken()
			return newAgn(v, expression(lex))
		} else {
			panic(fmt.Errorf("invalid assign form"))
		}
	}
	return e
}

// TopLevel
// 入力 - 評価 - 表示
func TopLevel(lex *Lex) (r bool) {
	r = false
	defer func() {
		err := recover()
		if err != nil {
			mes, ok := err.(string)
			if ok && mes == "quit" {
				r = true
			} else {
				fmt.Fprintln(os.Stderr, err)
				for lex.Token != ';' {
					lex.getToken()
				}
			}
		}
	}()
	for {
		fmt.Print("Calc> ")
		lex.getToken()
		if lex.Token == DEF {
			defineFunc(lex)
		} else {
			e := expression(lex)
			if lex.Token != ';' {
				panic(fmt.Errorf("invalid expression"))
			} else {
				fmt.Println(e.Eval(nil))
			}
		}
	}
	return r
}
