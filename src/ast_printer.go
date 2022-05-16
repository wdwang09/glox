package glox

import "fmt"

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return new(AstPrinter)
}

func (s *AstPrinter) Print(expr Expr) string {
	return expr.Accept(s).(string)
}

func (s *AstPrinter) parenthesize(name string, exprList ...Expr) string {
	res := "("
	res += name
	for _, expr := range exprList {
		res += " "
		res = res + expr.Accept(s).(string)
	}
	res += ")"
	return res
}

func (s *AstPrinter) visitAssignExpr(assign *Assign) interface{} {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitBinaryExpr(binary *Binary) interface{} {
	return s.parenthesize(binary.operator.lexeme, binary.left, binary.right)
}

func (s *AstPrinter) visitCallExpr(call *Call) interface{} {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitGetExpr(get *Get) interface{} {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitGroupingExpr(grouping *Grouping) interface{} {
	return s.parenthesize("group", grouping.expression)
}

func (s *AstPrinter) visitLiteralExpr(literal *Literal) interface{} {
	if literal.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", literal.value)
}

func (s *AstPrinter) visitLogicalExpr(logical *Logical) interface{} {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitSetExpr(set *Set) interface{} {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitSuperExpr(super *Super) interface{} {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitThisExpr(this *This) interface{} {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitUnaryExpr(unary *Unary) interface{} {
	return s.parenthesize(unary.operator.lexeme, unary.right)
}

func (s *AstPrinter) visitVariableExpr(variable *Variable) interface{} {
	// TODO implement me
	panic("implement me")
}
