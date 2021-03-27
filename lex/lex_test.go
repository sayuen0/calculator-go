package lex

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestLex_getToken(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want []interface{}
	}{
		// TODO: Add test cases.
		{
			name: "num + num",
			expr: "1 + 3 ;",
			want: []interface{}{Value(1), '+', Value(3), ';'},
		},
		{
			name: "num - num",
			expr: "1 - 3 ;",
			want: []interface{}{Value(1), '-', Value(3), ';'},
		},
		{
			name: "num * num",
			expr: "1 * 3 ;",
			want: []interface{}{Value(1), '*', Value(3), ';'},
		},
		{
			name: "num / num",
			expr: "1 / 3 ;",
			want: []interface{}{Value(1), '/', Value(3), ';'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ダミーファイル作成
			tmp, _ := ioutil.TempFile("", "tmp")
			orgStdin := os.Stdin
			os.Stdin = tmp
			defer func() {
				os.Remove(tmp.Name())
				os.Stdin = orgStdin
			}()
			tmp.Write([]byte(tt.expr))
			tmp.Seek(0, 0)
			var l Lex
			l.Init(tmp)
			// assertion scan したときにTokenの内容がセットされていること
			for i := 0; i < len(tt.want); i++ {
				l.getToken()
				switch a := tt.want[i].(type) {
				case Value:
					var n float64
					fmt.Sscan(l.TokenText(), &n)
					v := Value(n)
					if tt.want[i] != v {
						t.Errorf("failed, want %v, get %v", tt.want[i], v)
					}
				case rune:
					if tt.want[i] != l.Token {
						t.Errorf("failed, want %v, get %v", tt.want[i], l.Token)
					}
				default:
					t.Errorf("failed, invalid type %v", a)
				}
				log.Println(tt.want[i], l.Token, l.TokenText())
			}
		})
	}
}

// 短絡演算子評価テスト
func TestOps_Eval(t *testing.T) {
	type fields struct {
		code  rune
		left  Expr
		right Expr
	}
	type args struct {
		env *Env
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Value
	}{
		{
			name: "truthy and truthy",
			fields: fields{
				code:  AND,
				left:  Value(1.0),
				right: Value(1.0),
			},
			args: args{
				env: &Env{},
			},
			want: Value(1.0),
		},
		{
			name: "truthy and falsy",
			fields: fields{
				code:  AND,
				left:  Value(1.0),
				right: Value(0.0),
			},
			args: args{
				env: &Env{},
			},
			want: Value(0.0),
		},
		{
			name: "falsy and truthy",
			fields: fields{
				code:  AND,
				left:  Value(0.0),
				right: Value(1.0),
			},
			args: args{
				env: &Env{},
			},
			want: Value(0.0),
		},
		{
			name: "falsy and falsy",
			fields: fields{
				code:  AND,
				left:  Value(0.0),
				right: Value(0.0),
			},
			args: args{
				env: &Env{},
			},
			want: Value(0.0),
		},
		{
			name: "truthy or truthy",
			fields: fields{
				code:  OR,
				left:  Value(1.0),
				right: Value(1.0),
			},
			args: args{
				env: &Env{},
			},
			want: Value(1.0),
		},
		{
			name: "truthy or falsy",
			fields: fields{
				code:  OR,
				left:  Value(1.0),
				right: Value(0.0),
			},
			args: args{
				env: &Env{},
			},
			want: Value(1.0),
		},
		{
			name: "falsy or truthy",
			fields: fields{
				code:  OR,
				left:  Value(0.0),
				right: Value(1.0),
			},
			args: args{
				env: &Env{},
			},
			want: Value(1.0),
		},
		{
			name: "falsy or falsy",
			fields: fields{
				code:  AND,
				left:  Value(0.0),
				right: Value(0.0),
			},
			args: args{
				env: &Env{},
			},
			want: Value(0.0),
		},
		{
			name: "invalid ops code",
			fields: fields{
				code:  IF,
				left:  Value(0.0),
				right: Value(0.0),
			},
			args: args{
				env: &Env{},
			},
			want: Value(0.0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ops{
				code:  tt.fields.code,
				left:  tt.fields.left,
				right: tt.fields.right,
			}
			defer func() {
				err := recover()
				fmt.Printf("%T", err)
				if err != nil && err.(error).Error() != "invalid Ops code" {
					t.Errorf("panic %v", err)
				}
			}()
			if got := e.Eval(tt.args.env); got != tt.want {
				t.Errorf("Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSel_Eval(t *testing.T) {
	type fields struct {
		testForm Expr
		thenForm Expr
		elseForm Expr
	}
	type args struct {
		env *Env
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Value
	}{
		// TODO: Add test cases.
		{
			name: "truth testForm",
			fields: fields{
				testForm: Value(1.0),
				thenForm: Value(5.0),
				elseForm: Value(7.0),
			},
			args: args{
				env: &Env{},
			},
			want: Value(5.0),
		},
		{
			name: "falsy testForm",
			fields: fields{
				testForm: Value(0.0),
				thenForm: Value(8.0),
				elseForm: Value(12.0),
			},
			args: args{
				env: &Env{},
			},
			want: Value(12.0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Sel{
				testForm: tt.fields.testForm,
				thenForm: tt.fields.thenForm,
				elseForm: tt.fields.elseForm,
			}
			if got := e.Eval(tt.args.env); got != tt.want {
				t.Errorf("Eval() = %v, want %v", got, tt.want)
			}
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

func Test_expr2(t *testing.T) {
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
			if got := expr2(tt.args.lex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expr2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_expr3(t *testing.T) {
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
			if got := expr3(tt.args.lex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expr3() = %v, want %v", got, tt.want)
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

func Test_initKeyTable(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_makeSel(t *testing.T) {
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
			if got := makeSel(tt.args.lex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeSel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newOps(t *testing.T) {
	type args struct {
		code  rune
		left  Expr
		right Expr
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
			if got := newOps(tt.args.code, tt.args.left, tt.args.right); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newOps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newSel(t *testing.T) {
	type args struct {
		testForm Expr
		thenForm Expr
		elseForm Expr
	}
	tests := []struct {
		name string
		args args
		want *Sel
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSel(tt.args.testForm, tt.args.thenForm, tt.args.elseForm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSel() = %v, want %v", got, tt.want)
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
