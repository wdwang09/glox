package glox

import "fmt"

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (s *AstPrinter) PrintExpression(expr Expr) (string, error) {
	str, err := expr.accept(s)
	if err != nil {
		return "", err
	}
	return str.(string), nil
}

func (s *AstPrinter) PrintStatement(stmt Stmt) (string, error) {
	str, err := stmt.accept(s)
	if err != nil {
		return "", err
	}
	return str.(string), nil
}

func (s *AstPrinter) parenthesize(name string, exprList ...Expr) (string, error) {
	res := "("
	res += name
	for _, expr := range exprList {
		res += " "
		str, err := expr.accept(s)
		if err != nil {
			return "", err
		}
		res = res + str.(string)
	}
	res += ")"
	return res, nil
}

func (s *AstPrinter) parenthesize2(name string, parts ...interface{}) (string, error) {
	res := "("
	res += name + " "
	str, err := s.transform(parts...)
	if err != nil {
		return "", err
	}
	res += str
	res += ")"
	return res, nil
}

func (s *AstPrinter) transform(parts ...interface{}) (string, error) {
	res := ""
	var str interface{}
	var err error
	for i, part := range parts {
		if p, ok := part.(Expr); ok {
			str, err = p.accept(s)
		} else if p, ok := part.(Stmt); ok {
			str, err = p.accept(s)
		} else if p, ok := part.(*Token); ok {
			str = p.lexeme
		} else if p, ok := part.(*[]Expr); ok {
			var args []interface{}
			for _, arg := range *p {
				args = append(args, arg)
			}
			str, err = s.transform(args...)
			if err != nil {
				return "", err
			}
		} else {
			str = fmt.Sprint(part)
		}
		if err != nil {
			return "", err
		}
		if i != 0 && str != "" {
			res += " "
		}
		res += str.(string)
	}
	return res, nil
}

func (s *AstPrinter) visitAssignExpr(expr *Assign) (interface{}, error) {
	return s.parenthesize2("=", expr.name.lexeme, expr.value)
}

func (s *AstPrinter) visitBinaryExpr(expr *Binary) (interface{}, error) {
	return s.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (s *AstPrinter) visitCallExpr(expr *Call) (interface{}, error) {
	return s.parenthesize2("call", expr.callee, expr.arguments)
}

func (s *AstPrinter) visitGetExpr(expr *Get) (str interface{}, err error) {
	return s.parenthesize2(".", expr.object, expr.name.lexeme)
}

func (s *AstPrinter) visitGroupingExpr(expr *Grouping) (str interface{}, err error) {
	return s.parenthesize("group", expr.expression)
}

func (s *AstPrinter) visitLiteralExpr(expr *Literal) (str interface{}, err error) {
	if expr.value == nil {
		return "nil", nil
	}
	if str, ok := expr.value.(string); ok {
		return "\"" + str + "\"", nil
	}
	return fmt.Sprint(expr.value), nil
}

func (s *AstPrinter) visitLogicalExpr(expr *Logical) (str interface{}, err error) {
	return s.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (s *AstPrinter) visitSetExpr(expr *Set) (str interface{}, err error) {
	return s.parenthesize2("=", expr.object, expr.name.lexeme, expr.value)
}

func (s *AstPrinter) visitSuperExpr(expr *Super) (str interface{}, err error) {
	return s.parenthesize2("super", expr.method)
}

func (s *AstPrinter) visitThisExpr(expr *This) (str interface{}, err error) {
	return "this", nil
}

func (s *AstPrinter) visitUnaryExpr(expr *Unary) (str interface{}, err error) {
	return s.parenthesize(expr.operator.lexeme, expr.right)
}

func (s *AstPrinter) visitVariableExpr(expr *Variable) (str interface{}, err error) {
	return expr.name.lexeme, nil
}

func (s *AstPrinter) visitBlockStmt(stmt *Block) (interface{}, error) {
	res := "(block "
	for i, statement := range *stmt.statements {
		if i != 0 {
			res += " "
		}
		str, err := statement.accept(s)
		if err != nil {
			return nil, err
		}
		res += str.(string)
	}
	res += ")"
	return res, nil
}

func (s *AstPrinter) visitClassStmt(stmt *Class) (interface{}, error) {
	res := "(class " + stmt.name.lexeme
	if stmt.superclass != nil {
		res += " < "
		str, err := s.PrintExpression(stmt.superclass)
		if err != nil {
			return "", err
		}
		res += str
	}
	for _, method := range *stmt.methods {
		res += " "
		str, err := s.PrintStatement(method)
		if err != nil {
			return "", err
		}
		res += str
	}
	res += ")"
	return res, nil
}

func (s *AstPrinter) visitExpressionStmt(stmt *Expression) (interface{}, error) {
	return s.parenthesize(";", stmt.expression)
}

func (s *AstPrinter) visitFunctionStmt(stmt *Function) (interface{}, error) {
	res := "(fun " + stmt.name.lexeme + " ("
	for _, param := range *stmt.params {
		if param != (*stmt.params)[0] {
			res += " "
		}
		res += param.lexeme
	}
	res += ") "
	for i, body := range *stmt.body {
		if i != 0 {
			res += " "
		}
		str, err := body.accept(s)
		if err != nil {
			return "", err
		}
		res += str.(string)
	}
	res += ")"
	return res, nil
}

func (s *AstPrinter) visitIfStmt(stmt *If) (interface{}, error) {
	if stmt.elseBranch == nil {
		return s.parenthesize2("if", stmt.condition, stmt.thenBranch)
	}
	return s.parenthesize2("if-else", stmt.condition, stmt.thenBranch, stmt.elseBranch)
}

func (s *AstPrinter) visitPrintStmt(stmt *Print) (interface{}, error) {
	return s.parenthesize("print", stmt.expression)
}

func (s *AstPrinter) visitReturnStmt(stmt *Return) (interface{}, error) {
	if stmt.value == nil {
		return "(return)", nil
	}
	return s.parenthesize("return", stmt.value)
}

func (s *AstPrinter) visitVarStmt(stmt *Var) (interface{}, error) {
	if stmt.initializer == nil {
		return s.parenthesize2("var", stmt.name)
	}
	return s.parenthesize2("var", stmt.name, "=", stmt.initializer)
}

func (s *AstPrinter) visitWhileStmt(stmt *While) (interface{}, error) {
	return s.parenthesize2("while", stmt.condition, stmt.body)
}
