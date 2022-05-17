package glox

import "fmt"

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return new(AstPrinter)
}

func (s *AstPrinter) Print(expr Expr) (string, error) {
	sInterface, err := expr.Accept(s)
	if err != nil {
		return "", err
	}
	return sInterface.(string), nil
}

func (s *AstPrinter) parenthesize(name string, exprList ...Expr) (string, error) {
	res := "("
	res += name
	for _, expr := range exprList {
		res += " "
		sInterface, err := expr.Accept(s)
		if err != nil {
			return "", err
		}
		res = res + sInterface.(string)
	}
	res += ")"
	return res, nil
}

func (s *AstPrinter) visitAssignExpr(assign *Assign) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitBinaryExpr(binary *Binary) (interface{}, error) {
	return s.parenthesize(binary.operator.lexeme, binary.left, binary.right)
}

func (s *AstPrinter) visitCallExpr(call *Call) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitGetExpr(get *Get) (str interface{}, err error) {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitGroupingExpr(grouping *Grouping) (str interface{}, err error) {
	return s.parenthesize("group", grouping.expression)
}

func (s *AstPrinter) visitLiteralExpr(literal *Literal) (str interface{}, err error) {
	if literal.value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", literal.value), nil
}

func (s *AstPrinter) visitLogicalExpr(logical *Logical) (str interface{}, err error) {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitSetExpr(set *Set) (str interface{}, err error) {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitSuperExpr(super *Super) (str interface{}, err error) {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitThisExpr(this *This) (str interface{}, err error) {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitUnaryExpr(unary *Unary) (str interface{}, err error) {
	return s.parenthesize(unary.operator.lexeme, unary.right)
}

func (s *AstPrinter) visitVariableExpr(variable *Variable) (str interface{}, err error) {
	// TODO implement me
	panic("implement me")
}
