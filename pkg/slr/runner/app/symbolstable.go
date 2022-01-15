package app

import (
	"fmt"
)

type Symbol struct {
	kind   string
	name   string
	length int
	level  int
}

type SymbolTable struct {
	table []Symbol
}

type Stack struct {
	stack []SymbolTable
}

func (t *SymbolTable) AddSymbol(s Symbol) { // типо добавление символа в блок (A2)
	fmt.Println("Add Symbol name = ", s.name, " type = ", s.kind, " level = ", s.level)
	t.table = append(t.table, s)
}

func (t *SymbolTable) Find(varName string) Symbol { // A3
	for _, n := range t.table {
		if varName == n.name {
			return n
		}
	}
	return Symbol{}
}

func (s *Stack) CreateStack() { // типо конструктор
	newTable := SymbolTable{}
	s.stack = append(s.stack, newTable)
}

func (s *Stack) AddTable() { // A1
	newTable := SymbolTable{}
	s.stack = append(s.stack, newTable)
}

func (s *Stack) DeleteLast() { // A4
	s.stack = s.stack[:len(s.stack)-1] //Удаление последнего элемента
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func ContainsSymbol(a []Symbol, x string) bool {
	for _, n := range a {
		if x == n.name {

			return true
		}
	}
	return false
}

func (s *Stack) IsNotExistVar(symName string) bool {
	for level := 0; level < len(s.stack); level++ {
		if ContainsSymbol(s.stack[level].table, symName) {
			return false
		}
	}
	return true
}

func main() {
	types := []string{"int", "float", "dbl", "string"}

	mainStack := Stack{}
	mainStack.CreateStack()
	input := []string{"int", "id", "{", "float", "a", ";", "a", "int", "c", "{", "int", "t", ";", "t", "}", "c", "c", "c", "}", "int", "b", "b", "b"}
	for in := 0; in < len(input); in++ {
		//fmt.Println(input[in]);

		if input[in] == "{" {
			mainStack.AddTable()
			continue
		}

		if input[in] == "}" {
			if len(mainStack.stack) != 1 {
				mainStack.DeleteLast()
				continue
			} else {
				fmt.Println("ERROR!!! TRY TO DELETE MAIN STACK!!!")
			}
		}

		if Contains(types, input[in]) {
			isNotExistVar := mainStack.IsNotExistVar(input[in+1])
			if (in != len(input)-1) && isNotExistVar {

				sym := Symbol{}
				sym.kind = input[in]
				sym.level = len(mainStack.stack) - 1

				if !Contains(types, input[in+1]) && input[in+1] != ";" {
					sym.name = input[in+1]
					in = in + 1

					mainStack.stack[len(mainStack.stack)-1].AddSymbol(sym)
					continue
				} else {
					fmt.Println("ERROR!!! EXPECTED ID BUT NULL GIVEN!!!")
					break
				}
			} else if !isNotExistVar {
				fmt.Println("ERROR!!! THIS ID ALREDY DECLARED!!! Name = ", input[in+1])

				break
			}
		} else if input[in] != ";" {
			if !mainStack.IsNotExistVar(input[in]) {

				requestedVar := mainStack.stack[len(mainStack.stack)-1].Find(input[in])

				fmt.Println("Requested Var name = ", requestedVar.name, " type = ", requestedVar.kind, " level = ", requestedVar.level)
				continue

			} else {
				fmt.Println("ERROR!!! UNKNOWN VAR!!! Name = ", input[in])
				break
			}
		}
	}

	fmt.Println(mainStack.stack)
}
