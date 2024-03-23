package parser

import (
	"fmt"
	"strconv"

	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

const (
	_           int = iota
	lowest          // 最低优先级
	equals          // ==
	lessGreater     // ！=
	sum             // +
	product         // *
	prefix          // -X or !X
	call            // function call
	index           // array index
)

// 优先级
var precedences = map[token.TypeToken]int{
	token.EQ:       equals,
	token.NOT_EQ:   equals,
	token.LT:       lessGreater,
	token.GT:       lessGreater,
	token.PLUS:     sum,
	token.MINUS:    sum,
	token.SLASH:    product,
	token.ASTERISK: product,
	token.LPAREN:   call,
	token.LBRACKET: index,
}

type (
	prefixParseFunc func() ast.Expression               // 前缀解析函数
	infixParseFunc  func(ast.Expression) ast.Expression // 中缀解析函数
)

// Parser 解析器
type Parser struct {
	l      *lexer.Lexer // 词法分析器
	errors []string     // 错误信息

	curToken  token.Token // 当前token
	peekToken token.Token // 下一个token

	prefixParseFns map[token.TypeToken]prefixParseFunc // 前缀解析函数
	infixParseFns  map[token.TypeToken]infixParseFunc  // 中缀解析函数
}

// New 创建解析器
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: make([]string, 0),
	}

	// 初始化当前和下一个token
	p.nextToken()
	p.nextToken()

	// 初始化prefixParseFns
	p.prefixParseFns = map[token.TypeToken]prefixParseFunc{}
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	// 初始化infixParseFns
	p.infixParseFns = map[token.TypeToken]infixParseFunc{}
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	return p
}

// Errors 获取错误信息
func (p *Parser) Errors() []string {
	return p.errors
}

// peekError 预读错误
func (p *Parser) peekError(t token.TypeToken) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// nextToken 获取下一个token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram 解析程序
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = make([]ast.Statement, 0)
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// parseStatement 解析语句
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement 解析let语句
func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()

	// 解析表达式
	stmt.Value = p.parseExpression(lowest)

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement 解析return语句
func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// 解析表达式
	stmt.ReturnValue = p.parseExpression(lowest)

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// curTokenIs 判断当前token是否为指定token
func (p *Parser) curTokenIs(t token.TypeToken) bool {
	return p.curToken.Type == t
}

// peekTokenIs 判断下一个token是否为指定token
func (p *Parser) peekTokenIs(t token.TypeToken) bool {
	return p.peekToken.Type == t
}

// expectPeek 判断下一token是否为指定token，如果是则获取下一个token
func (p *Parser) expectPeek(t token.TypeToken) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// registerPrefix 注册前缀解析函数
func (p *Parser) registerPrefix(tokenType token.TypeToken, fn prefixParseFunc) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix 注册中缀解析函数
func (p *Parser) registerInfix(tokenType token.TypeToken, fn infixParseFunc) {
	p.infixParseFns[tokenType] = fn
}

// parseExpressionStatement 解析表达式语句
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(lowest)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// parseExpression 解析表达式
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// parseIdentifier 解析标识符
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parseIntegerLiteral 解析整数字面量
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

// parseBoolean 解析布尔值
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// parseGroupedExpression 解析括号表达式
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(lowest)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

// parsePrefixExpression 解析前缀表达式
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()
	expression.Right = p.parseExpression(prefix)
	return expression
}

// parseInfixExpression 解析中缀表达式
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{Token: p.curToken, Left: left, Operator: p.curToken.Literal}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

// parseIfExpression 解析if表达式
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(lowest)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()
	if !p.curTokenIs(token.RBRACE) {
		return nil
	}
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}
	if !p.curTokenIs(token.RBRACE) {
		return nil
	}
	return expression
}

// parseBlockStatement 解析块语句
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}
	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

// parseFunctionLiteral 解析函数表达式
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

// parseFunctionParameters 解析函数参数
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}
	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return identifiers
}

// parseCallExpression 解析调用表达式
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

// parseIndexExpression 解析索引表达式
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}
	p.nextToken()
	exp.Index = p.parseExpression(lowest)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return exp
}

// parseStringLiteral 解析字符串字面量
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parseArrayLiteral 解析数组字面量
func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

// parseExpressionList 解析表达式列表
func (p *Parser) parseExpressionList(end token.TypeToken) []ast.Expression {
	var list []ast.Expression
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(lowest))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(lowest))
	}
	if !p.expectPeek(end) {
		return nil
	}
	return list
}

// noPrefixParseFnError 未找到前缀解析函数
func (p *Parser) noPrefixParseFnError(t token.TypeToken) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// peekPrecedence 获取下一个token的优先级
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return lowest
}

// curPrecedence 获取当前token的优先级
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return lowest
}
