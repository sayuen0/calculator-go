package lex

func boolToValue(x bool) Value {
	if x {
		return 1.0
	} else {
		return 0.0
	}
}

func isTrue(x Value) bool{
	return x != 0.0
}

func isFalse(x Value) bool{
	return x == 0.0
}
