package app

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

func (t *SymbolTable) AddSymbol(s Symbol) {
	t.table = append(t.table, s)
}

func (t *SymbolTable) Find(varName string) *Symbol {
	for _, n := range t.table {
		if varName == n.name {
			return &n
		}
	}
	return nil
}

func (s *Stack) CreateStack() {
	newTable := SymbolTable{}
	s.stack = append(s.stack, newTable)
}

func (s *Stack) GetLast() *SymbolTable {
	return &s.stack[len(s.stack)-1]
}

func (s *Stack) AddTable() {
	newTable := SymbolTable{}
	s.stack = append(s.stack, newTable)
}

func (s *Stack) DeleteLast() {
	s.stack = s.stack[:len(s.stack)-1]
}

func (s *Stack) IsNotExistVar(symName string) bool {
	for level := 0; level < len(s.stack); level++ {
		if ContainsSymbol(s.stack[level].table, symName) {
			return false
		}
	}
	return true
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
