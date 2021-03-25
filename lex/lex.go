package lex

import (
	"fmt"
	"os"
	"text/scanner"
)

// Lexとは、レキシカルアナライザ(字句解析プログラム)ジェネレータのこと
type Lex struct {
	scanner.Scanner
	Token rune
}

// 標準入力を1つ読み込んでruneを持つ
func (s *Lex) getToken() {
	s.Token = s.Scan()
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
		fmt.Print("Calc>")
		lex.getToken()
		val := expression(lex)
		if lex.Token != ';' {
			panic(fmt.Errorf("invalid expression"))
		} else {
			fmt.Println(val)
		}
	}
	return r
}

// expression: 式
// 式  = 項 { ("+" | "-"), 項 }.
func expression(lex *Lex) float64 {
	val := term(lex)
	// + か - 以外が現れるまでループ
	for {
		switch lex.Token {
		case '+':
			lex.getToken()
			val += term(lex)
		case '-':
			lex.getToken()
			val -= term(lex)
		default:
			return val
		}
	}
}

// term: 項
// 項  = 因子 { ("*" | "/"), 因子 }.
func term(lex *Lex) float64 {
	val := factor(lex)
	for {
		switch lex.Token {
		case '*':
			lex.getToken()
			val *= factor(lex)
		case '/':
			lex.getToken()
			val /= factor(lex)
		default:
			return val
		}
	}
}

// factor: 因子
// 因子 = 数値 | ("+" | "-"), 因子 | "(" 式 ")".
func factor(lex *Lex) float64 {
	switch lex.Token {
	case '(':
		lex.getToken()
		val := expression(lex)
		if lex.Token != ')' {
			panic(fmt.Errorf("'(' expected"))
		}
		lex.getToken()
		return val
	case '+':
		lex.getToken()
		return factor(lex)
	case '-':
		lex.getToken()
		return -factor(lex)
	case scanner.Int, scanner.Float:
		var n float64
		fmt.Sscan(lex.TokenText(), &n)
		lex.getToken()
		return n
	case scanner.Ident:
		text := lex.TokenText()
		if text == "quit" {
			panic(text)
		}
		fallthrough
	default:
		panic(fmt.Errorf("unexpected token: %v", lex.TokenText()))
	}
}
