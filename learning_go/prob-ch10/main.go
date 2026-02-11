package main

import "golang.org/x/exp/constraints"

type Number interface {
	constraints.Integer | constraints.Float
}

// Add 関数は、2つの整数 a と b を受け取り、それらの合計を返します。
// より詳しい情報については、次のリンクを参照してください:
// https://www.mathsisfun.com/numbers/addition.html
func Add[T Number](a, b T) T {
	return a + b
}
