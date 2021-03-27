package lex

import (
	"fmt"
	"log"
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

// キーワード
const (
	DEF = -(iota + 10)
	END
	IF
	THEN
	ELSE
	NOT
	AND
	OR
	EQ
	NE
	LT
	GT
	LE
	GE
	BGN
	WHL
	DO
	LET
	IN
)

var keyTable = make(map[string]rune)

func initKeyTable() {
	keyTable["def"] = DEF
	keyTable["end"] = END
	keyTable["if"] = IF
	keyTable["then"] = THEN
	keyTable["else"] = ELSE
	keyTable["and"] = AND
	keyTable["or"] = OR
	keyTable["not"] = NOT
	keyTable["begin"] = BGN
	keyTable["while"] = WHL
	keyTable["do"] = DO
	keyTable["let"] = LET
	keyTable["in"] = IN
}

// 短絡演算子
type Ops struct {
	code        rune
	left, right Expr
}

func newOps(code rune, left, right Expr) Expr {
	return &Ops{code, left, right}
}

// 短絡演算子の評価
func (e *Ops) Eval(env *Env) Value {
	x := e.left.Eval(env)
	switch e.code {
	case AND:
		if isTrue(x) {
			return e.right.Eval(env)
		}
		return x
	case OR:
		if isTrue(x) {
			return x
		}
		return e.right.Eval(env)
	default:
		panic(fmt.Errorf("invalid Ops code"))
	}
}

// if
type Sel struct {
	testForm, thenForm, elseForm Expr
}

func newSel(testForm, thenForm, elseForm Expr) *Sel {
	return &Sel{testForm, thenForm, elseForm}
}

// if式の評価
func (e *Sel) Eval(env *Env) Value {
	if isTrue(e.testForm.Eval(env)) {
		return e.thenForm.Eval(env)
	}
	return e.elseForm.Eval(env)
}

// ifの処理
func makeSel(lex *Lex) Expr {
	testForm := expression(lex)
	if lex.Token == THEN {
		lex.getToken()
		thenForm := expression(lex)
		switch lex.Token {
		case ELSE:
			lex.getToken()
			elseForm := expression(lex)
			if lex.Token != END {
				panic(fmt.Errorf("'end' expected"))
			}
			lex.getToken()
			return newSel(testForm, thenForm, elseForm)
		case END:
			lex.getToken()
			return newSel(testForm, thenForm, Value(0.0))
		default:
			panic(fmt.Errorf("'else' or 'end' expected"))
		}
	} else {
		panic(fmt.Errorf("'then' expected"))
	}
	// このブロックには到達しない
	return nil
}

// 標準入力を1つ読み込んでruneを持つ
func (lex *Lex) getToken() {
	lex.Token = lex.Scan()
	switch lex.Token {
	case scanner.Ident:
		key, ok := keyTable[lex.TokenText()]
		if ok {
			lex.Token = key
		}
	case '=':
		if lex.Peek() == '=' {
			lex.Next()
			lex.Token = EQ
		}
	case '!':
		if lex.Peek() == '=' {
			lex.Next()
			lex.Token = NE
		} else {
			lex.Token = NOT
		}
	case '<':
		if lex.Peek() == '=' {
			lex.Next()
			lex.Token = LE
		} else {
			lex.Token = LT
		}
	case '>':
		if lex.Peek() == '=' {
			lex.Next()
			lex.Token = GE
		} else {
			lex.Token = GT
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
	case NOT:
		lex.getToken()
		return newOp1(NOT, factor(lex))
	case IF:
		lex.getToken()
		return makeSel(lex)
	case BGN:
		lex.getToken()
		return makeBegin(lex)
	case WHL:
		lex.getToken()
		return makeWhile(lex)
	case LET:
		lex.getToken()
		return makeLet(lex)
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

// 論理演算子
func expr1(lex *Lex) Expr {
	e := expr2(lex)
	for {
		x := lex.Token
		switch x {
		case AND, OR:
			lex.getToken()
			e = newOps(x, e, expr2(lex))
		default:
			return e
		}
	}
}

// 比較演算子
func expr2(lex *Lex) Expr {
	e := expr3(lex)
	x := lex.Token
	switch x {
	case EQ, NE, LT, GT, LE, GE:
		lex.getToken()
		return newOp2(x, e, expr3(lex))
	default:
		return e
	}
}

// 式
func expr3(lex *Lex) Expr {
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
				for {
					c:= lex.Peek()
					if c == '\n' {break }
					lex.Next()
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
				log.Println(lex.TokenText())
				panic(fmt.Errorf("invalid expression"))
			} else {
				fmt.Println(e.Eval(nil))
			}
		}
	}
	return r
}
