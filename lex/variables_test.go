package lex

import (
	"math"
	"reflect"
	"runtime"
	"testing"
)

func TestAgn_Eval(t *testing.T) {
	type fields struct {
		name Variable
		expr Expr
	}
	tests := []struct {
		name   string
		fields fields
		want   Value
	}{
		// TODO: Add test cases.
		{
			name: "case1",
			fields: fields{
				name: Variable("a"),
				expr: newOp1('+', Value(1)),
			},
			want: Value(1),
		},
		{
			name: "case2",
			fields: fields{
				name: Variable("b"),
				expr: newOp1('-', Value(1)),
			},
			want: Value(-1),
		},
		{
			name: "case3",
			fields: fields{
				name: Variable("c"),
				expr: newOp1('+', Variable("a")),
			},
			want: Value(1),
		},
		{
			name: "case4",
			fields: fields{
				name: Variable("d"),
				expr: newOp2('*', Variable("a"), Variable("b")),
			},
			want: Value(-1),
		},
		{
			name: "case5",
			fields: fields{
				name: Variable("e"),
				expr: newOp2('+', Value(0), Value(334)),
			},
			want: Value(114514),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Agn{
				name: tt.fields.name,
				expr: tt.fields.expr,
			}
			if got := a.Eval(); got != tt.want {
				t.Errorf("Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_Eval(t *testing.T) {
	type fields struct {
		fn Func
		xs []Expr
	}
	tests := []struct {
		name   string
		fields fields
		want   Value
	}{
		// TODO: Add test cases.
		{
			name: "case1",
			fields: fields{
				fn: Func1(math.Sqrt),
				xs: []Expr{Value(4)},
			},
			want: Value(math.Sqrt(4)),
		},
		{
			name: "case2",
			fields: fields{
				fn: Func1(math.Sin),
				xs: []Expr{Value(math.Pi / 2)},
			},
			want: Value(math.Sin(math.Pi / 2)),
		},
		{
			name: "case3",
			fields: fields{
				fn: Func2(math.Pow),
				xs: []Expr{Value(2), Value(4)},
			},
			want: Value(math.Pow(2, 4)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				fn: tt.fields.fn,
				xs: tt.fields.xs,
			}
			if got := a.Eval(); got != tt.want {
				t.Errorf("Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunc1_Argc(t *testing.T) {
	tests := []struct {
		name string
		f    Func1
		want int
	}{
		{
			name: "case",
			f:    Func1(math.Sin),
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Argc(); got != tt.want {
				t.Errorf("Argc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunc2_Argc(t *testing.T) {
	tests := []struct {
		name string
		f    Func2
		want int
	}{
		{
			name: "case",
			f:    Func2(math.Pow),
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Argc(); got != tt.want {
				t.Errorf("Argc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitFunc(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "sqrt",
			want: "math.Sqrt",
		},
		{
			name: "sin",
			want: "math.Sin",
		},
		{
			name: "cos",
			want: "math.Cos",
		},
		{
			name: "tan",
			want: "math.Tan",
		},
		{
			name: "sinh",
			want: "math.Sinh",
		},
		{
			name: "cosh",
			want: "math.Cosh",
		},
		{
			name: "tanh",
			want: "math.Tanh",
		},
		{
			name: "asin",
			want: "math.Asin",
		},
		{
			name: "acos",
			want: "math.Acos",
		},
		{
			name: "atan",
			want: "math.Atan",
		},
		{
			name: "atan2",
			want: "math.Atan2",
		},
		{
			name: "exp",
			want: "math.Exp",
		},
		{
			name: "pow",
			want: "math.Pow",
		},
		{
			name: "log",
			want: "math.Log",
		},
		{
			name: "log10",
			want: "math.Log10",
		},
		{
			name: "log2",
			want: "math.Log2",
		},
	}
	InitFunc()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runtime.FuncForPC(reflect.ValueOf(funcTable[tt.name]).Pointer()).Name()
			if got != tt.want {
				t.Errorf("InitFunc() failed, got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariable_Eval(t *testing.T) {
	tests := []struct {
		name string
		v    Variable
		want Value
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 事前に変数を代入していないといけない
			// やらなくてもいいけど、代入してないならpanicするといいよね
			if got := tt.v.Eval(); got != tt.want {
				t.Errorf("Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newAgn(t *testing.T) {
	type args struct {
		v Variable
		e Expr
	}
	tests := []struct {
		name string
		args args
		want *Agn
	}{
		{
			name: "case",
			args: args{
				v: Variable("a"),
				e: Value(1),
			},
			want: &Agn{
				name: Variable("a"),
				expr: Value(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newAgn(tt.args.v, tt.args.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newAgn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newApp(t *testing.T) {
	type args struct {
		fn Func
		xs []Expr
	}
	tests := []struct {
		name string
		args args
		want *App
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newApp(tt.args.fn, tt.args.xs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newApp() = %v, want %v", got, tt.want)
			}
		})
	}
}
