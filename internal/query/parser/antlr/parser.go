package antlrparser

import (
	"context"
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"streaming-golang/internal/app/apperr"
	"streaming-golang/internal/app/transactional"
	"streaming-golang/internal/domain"
	"streaming-golang/internal/query/parser/antlr/generated"
)

type Parser struct{}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(ctx context.Context, expressions []string) (transactional.FilterSet, error) {
	nodes := make([]domain.FilterNode, 0, len(expressions))
	for _, expression := range expressions {
		if err := ctx.Err(); err != nil {
			return transactional.FilterSet{}, err
		}
		parsed, err := p.parseExpression(expression)
		if err != nil {
			return transactional.FilterSet{}, err
		}
		nodes = append(nodes, parsed...)
	}

	return transactional.FilterSet{Expressions: expressions, Nodes: nodes}, nil
}

func (p *Parser) parseExpression(expression string) ([]domain.FilterNode, error) {
	if strings.TrimSpace(expression) == "" {
		return nil, apperr.New(apperr.Invalid, "filter expression cannot be empty")
	}

	input := antlr.NewInputStream(expression)
	lexer := generated.NewOutboundAPILexer(input)
	errors := &syntaxErrorListener{expression: expression}

	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errors)

	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := generated.NewOutboundAPIParser(tokens)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(errors)

	tree := parser.ExpressionsSection()
	if errors.hasErrors() {
		return nil, apperr.New(apperr.Invalid, errors.message())
	}
	if parser.HasError() {
		return nil, apperr.New(apperr.Invalid, "invalid filter expression")
	}

	result := tree.Accept(newASTVisitor())
	nodes, ok := result.([]domain.FilterNode)
	if !ok {
		return nil, apperr.New(apperr.Invalid, "invalid filter expression")
	}
	return nodes, nil
}

type syntaxErrorListener struct {
	*antlr.DefaultErrorListener
	expression string
	errors     []syntaxError
}

type syntaxError struct {
	line   int
	column int
	text   string
}

func (l *syntaxErrorListener) SyntaxError(_ antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, _ antlr.RecognitionException) {
	l.errors = append(l.errors, syntaxError{
		line:   line,
		column: column,
		text:   fmt.Sprintf("%s near %q", msg, offendingText(offendingSymbol)),
	})
}

func (l *syntaxErrorListener) hasErrors() bool {
	return len(l.errors) > 0
}

func (l *syntaxErrorListener) message() string {
	if len(l.errors) == 0 {
		return ""
	}

	first := l.errors[0]
	return fmt.Sprintf("invalid filter expression at line %d, column %d: %s", first.line, first.column, first.text)
}

func offendingText(symbol interface{}) string {
	token, ok := symbol.(antlr.Token)
	if !ok || token == nil {
		return ""
	}
	return token.GetText()
}
