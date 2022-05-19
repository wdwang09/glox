package glox

import "fmt"

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
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

func (s *AstPrinter) visitAssignExpr(expr *Assign) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitBinaryExpr(expr *Binary) (interface{}, error) {
	return s.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (s *AstPrinter) visitCallExpr(expr *Call) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitGetExpr(expr *Get) (str interface{}, err error) {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitGroupingExpr(expr *Grouping) (str interface{}, err error) {
	return s.parenthesize("group", expr.expression)
}

func (s *AstPrinter) visitLiteralExpr(expr *Literal) (str interface{}, err error) {
	if expr.value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.value), nil
}

func (s *AstPrinter) visitLogicalExpr(expr *Logical) (str interface{}, err error) {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitSetExpr(expr *Set) (str interface{}, err error) {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitSuperExpr(expr *Super) (str interface{}, err error) {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitThisExpr(expr *This) (str interface{}, err error) {
	// TODO implement me
	panic("implement me")
}

func (s *AstPrinter) visitUnaryExpr(expr *Unary) (str interface{}, err error) {
	return s.parenthesize(expr.operator.lexeme, expr.right)
}

func (s *AstPrinter) visitVariableExpr(expr *Variable) (str interface{}, err error) {
	// TODO implement me
	panic("implement me")
}
