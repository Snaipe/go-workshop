package main

import (
	"fmt"
	"strconv"
)

func Eval(tokens []Token) ([]int, error) {

	var results []int
	for len(tokens) > 0 {
		t := tokens[0]
		tokens = tokens[1:]

		switch t.Type {
		case TokenNumber:
			n, err := strconv.Atoi(t.Val)
			if err != nil {
				return nil, err
			}
			results = append(results, n)
		case TokenSymbol:
			return nil, fmt.Errorf("unexpected symbol %q", t.Val)
		case TokenLParen:
			rparen := -1
			depth := 1
			for i := 1; i < len(tokens); i++ {
				switch tokens[i].Type {
				case TokenLParen:
					depth++
				case TokenRParen:
					depth--
				}
				if depth == 0 {
					rparen = i
					break
				}
			}
			if rparen == -1 {
				return nil, fmt.Errorf("mismatching parenthesis; got end of string but expected ')'")
			}
			ltokens := tokens[:rparen]

			if len(ltokens) < 2 {
				return nil, fmt.Errorf("expected at least 1 argument in function call")
			}

			sym := ltokens[0]
			if sym.Type != TokenSymbol {
				return nil, fmt.Errorf("expected a symbol, got %q", sym.Val)
			}

			var fn func([]int) (int, error)
			switch sym.Val {
			case "+":
				fn = sum
			case "-":
				fn = sub
			case "*":
				fn = mul
			case "/":
				fn = div
			default:
				return nil, fmt.Errorf("unknown function %q", sym.Val)
			}

			args, err := Eval(ltokens[1:])
			if err != nil {
				return nil, err
			}

			n, err := fn(args)
			if err != nil {
				return nil, err
			}
			results = append(results, n)

			tokens = tokens[rparen+1:]
		default:
			return nil, fmt.Errorf("unexpected token %q", t.Val)
		}
	}
	return results, nil
}

func sum(args []int) (int, error) {
	if len(args) == 0 {
		return 0, fmt.Errorf("not enough arguments for '+'")
	}
	var sum int
	for _, n := range args {
		sum += n
	}
	return sum, nil
}

func mul(args []int) (int, error) {
	if len(args) == 0 {
		return 0, fmt.Errorf("not enough arguments for '*'")
	}
	product := 1
	for _, n := range args {
		product *= n
	}
	return product, nil
}

func sub(args []int) (int, error) {
	switch len(args) {
	case 0:
		return 0, fmt.Errorf("not enough arguments for '-'")
	case 1:
		return -args[0], nil
	default:
		result := args[0]
		for _, n := range args[1:] {
			result -= n
		}
		return result, nil
	}
}

func div(args []int) (int, error) {
	switch len(args) {
	case 0:
		return 0, fmt.Errorf("not enough arguments for '/'")
	case 1:
		return -args[0], nil
	default:
		num := args[0]
		denom := 1
		for _, n := range args[1:] {
			denom *= n
		}
		return num / denom, nil
	}
}
