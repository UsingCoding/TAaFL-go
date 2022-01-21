package common

type ASTRule string

const (
	VariableDeclaration ASTRule = "variable-declaration"
	VariableOperand     ASTRule = "variable-operand"
	AssigmentExpression ASTRule = "assigment-expression"
	UnaryExpression     ASTRule = "unary-expression"
	Expression          ASTRule = "expression" // for instance a = 4 expression

	BeginBlockStatement ASTRule = "begin-block-statement"
	BlockStatement      ASTRule = "block-statement"

	Condition   ASTRule = "condition"
	IFStatement ASTRule = "if-statement"

	UpdateExpression ASTRule = "update-expression"
	FORStatement     ASTRule = "for-statement"

	Addition       ASTRule = "addition"
	Subtraction    ASTRule = "subtraction"
	Multiplication ASTRule = "multiplication"
	Division       ASTRule = "division"
)
