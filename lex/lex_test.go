package lex

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestLex_getToken(t *testing.T) {
	tests := []struct {
		name   string
		expr   string
		want   []rune
	}{
		// TODO: Add test cases.
		{name: "case1",
		expr: "1 + 3 ;",
		want: []rune{'1', '+', '3', ';'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ダミーファイル作成
			tmp, _ := ioutil.TempFile("","tmp")
			orgStdin := os.Stdin
			os.Stdin = tmp
			defer func() {
				os.Remove(tmp.Name())
				os.Stdin = orgStdin
			}()
			tmp.Write([]byte(tt.expr))
			tmp.Seek(0,0)
			var l Lex
			l.Init(tmp)
			// assertion scan したときにTokenの内容がセットされていること
			// TODO: assertion keytableの内容を加味
			l.getToken()
			log.Println(l.Token, l.TokenText())
			l.getToken()
			log.Println(l.Token, l.TokenText())
			l.getToken()
			log.Println(l.Token, l.TokenText())
			l.getToken()
			log.Println(l.Token, l.TokenText())
		})
	}
}

func TestTopLevel(t *testing.T) {
	type args struct {
		lex *Lex
	}
	tests := []struct {
		name  string
		args  args
		wantR bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := TopLevel(tt.args.lex); gotR != tt.wantR {
				t.Errorf("TopLevel() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func Test_expr1(t *testing.T) {
	type args struct {
		lex *Lex
	}
	tests := []struct {
		name string
		args args
		want Expr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := expr1(tt.args.lex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expr1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_expression(t *testing.T) {
	type args struct {
		lex *Lex
	}
	tests := []struct {
		name string
		args args
		want Expr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := expression(tt.args.lex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expression() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_factor(t *testing.T) {
	type args struct {
		lex *Lex
	}
	tests := []struct {
		name string
		args args
		want Expr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := factor(tt.args.lex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("factor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getArgs(t *testing.T) {
	type args struct {
		lex *Lex
	}
	tests := []struct {
		name string
		args args
		want []Expr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getArgs(tt.args.lex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_term(t *testing.T) {
	type args struct {
		lex *Lex
	}
	tests := []struct {
		name string
		args args
		want Expr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := term(tt.args.lex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("term() = %v, want %v", got, tt.want)
			}
		})
	}
}
