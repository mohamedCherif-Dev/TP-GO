package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage : calc <nombre> <opérateur> <nombre>")
		fmt.Println("Opérateurs : + - x /")
		os.Exit(1)
	}

	num1, err1 := strconv.ParseFloat(os.Args[1], 64)
	num2, err2 := strconv.ParseFloat(os.Args[3], 64)

	if err1 != nil || err2 != nil {
		fmt.Println("Erreur : les opérandes doivent être des nombres valides")
		os.Exit(1)
	}

	op := os.Args[2]
	var resultat float64

	switch op {
	case "+":
		resultat = num1 + num2
	case "-":
		resultat = num1 - num2
	case "x", "*":
		resultat = num1 * num2
	case "/":
		if num2 == 0 {
			fmt.Println("Erreur : division par zéro impossible")
			os.Exit(1)
		}
		resultat = num1 / num2
	default:
		fmt.Println("Erreur : opérateur inconnu. Utilisez +, -, x ou /")
		os.Exit(1)
	}

	fmt.Printf("%g %s %g = %g\n", num1, op, num2, resultat)
}
