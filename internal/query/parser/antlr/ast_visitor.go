package antlrparser

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"streaming-golang/internal/domain"
	"streaming-golang/internal/query/parser/antlr/generated"
)

type astVisitor struct {
	*generated.BaseOutboundAPIParserVisitor
}

func newASTVisitor() *astVisitor {
	return &astVisitor{
		BaseOutboundAPIParserVisitor: &generated.BaseOutboundAPIParserVisitor{
			BaseParseTreeVisitor: &antlr.BaseParseTreeVisitor{},
		},
	}
}

func (v *astVisitor) VisitExpressionsSection(ctx *generated.ExpressionsSectionContext) interface{} {
	nodes := make([]domain.FilterNode, 0)
	for _, section := range ctx.AllKeyFilterSection() {
		nodes = appendFilterNodes(nodes, section.Accept(v))
	}
	return nodes
}

func (v *astVisitor) VisitKeyFilterSection(ctx *generated.KeyFilterSectionContext) interface{} {
	nodes := make([]domain.FilterNode, 0, len(ctx.AllKeyComparison()))
	for _, comparison := range ctx.AllKeyComparison() {
		nodes = appendFilterNodes(nodes, comparison.Accept(v))
	}
	return nodes
}

func (v *astVisitor) VisitIdPointInTimeArithmeticComparison(ctx *generated.IdPointInTimeArithmeticComparisonContext) interface{} {
	return domain.ComparisonFilter{
		Raw:      ctx.GetText(),
		Field:    fieldName(ctx.ID(), ctx.KeySurfaceColumn()),
		Operator: terminalText(ctx.COMPARISON_OPERATOR()),
		Value:    pointInTimeValue(ctx.PointInTimeArithmetic()),
	}
}

func (v *astVisitor) VisitIdTimeIntervalIn(ctx *generated.IdTimeIntervalInContext) interface{} {
	return domain.ComparisonFilter{
		Raw:      ctx.GetText(),
		Field:    fieldName(ctx.ID(), ctx.KeySurfaceColumn()),
		Operator: terminalText(ctx.IN()),
		Value:    timeIntervalValue(ctx.TimeIntervalArithmetic()),
	}
}

func (v *astVisitor) VisitIdNumericComparison(ctx *generated.IdNumericComparisonContext) interface{} {
	return domain.ComparisonFilter{
		Raw:      ctx.GetText(),
		Field:    terminalText(ctx.ID()),
		Operator: terminalText(ctx.COMPARISON_OPERATOR()),
		Value: domain.FilterValue{
			Kind: domain.FilterValueNumber,
			Raw:  tokenText(ctx.GetNumber()),
		},
	}
}

func (v *astVisitor) VisitIdLatestGlobalComparison(ctx *generated.IdLatestGlobalComparisonContext) interface{} {
	raw := nodeText(ctx.LatestGlobalFunction())
	return domain.ComparisonFilter{
		Raw:      ctx.GetText(),
		Field:    terminalText(ctx.ID()),
		Operator: terminalText(ctx.COMPARISON_OPERATOR()),
		Value: domain.FilterValue{
			Kind:     domain.FilterValueLatestGlobal,
			Raw:      raw,
			Function: parseFunctionName(raw),
		},
	}
}

func (v *astVisitor) VisitIdTimeIntervalToPointInTimeComparison(ctx *generated.IdTimeIntervalToPointInTimeComparisonContext) interface{} {
	return domain.ComparisonFilter{
		Raw:      ctx.GetText(),
		Field:    fieldName(ctx.ID(), ctx.KeySurfaceColumn()),
		Operator: terminalText(ctx.COMPARISON_OPERATOR()),
		Value:    intervalToPointInTimeValue(ctx.TimeIntervalToPointInTime()),
	}
}

func (v *astVisitor) VisitIdLatestComparison(ctx *generated.IdLatestComparisonContext) interface{} {
	return domain.ComparisonFilter{
		Raw:      ctx.GetText(),
		Field:    terminalText(ctx.ID()),
		Operator: terminalText(ctx.COMPARISON_OPERATOR()),
		Value:    latestValue(ctx.LatestFunction()),
	}
}

func (v *astVisitor) VisitTextComparison(ctx *generated.TextComparisonContext) interface{} {
	return domain.ComparisonFilter{
		Raw:      ctx.GetText(),
		Field:    nodeText(ctx.TextColumn()),
		Operator: terminalText(ctx.COMPARISON_OPERATOR()),
		Value: domain.FilterValue{
			Kind: domain.FilterValueText,
			Raw:  nodeText(ctx.GenericValue()),
		},
	}
}

func (v *astVisitor) VisitRankOver(ctx *generated.RankOverContext) interface{} {
	rankOver := ctx.RankOverFunction()
	filter := domain.RankOverFilter{Raw: nodeText(rankOver)}
	if rankOver == nil {
		return filter
	}

	ids := terminalTexts(rankOver.AllID())
	sortOrders := terminalTexts(rankOver.AllSORT_ORDER())
	partitionCount := len(ids) - len(sortOrders)
	if partitionCount < 0 {
		partitionCount = 0
	}
	filter.PartitionBy = append(filter.PartitionBy, ids[:partitionCount]...)
	for i, orderID := range ids[partitionCount:] {
		direction := ""
		if i < len(sortOrders) {
			direction = sortOrders[i]
		}
		filter.OrderBy = append(filter.OrderBy, domain.SortExpression{
			Field:     orderID,
			Direction: direction,
		})
	}
	for _, bound := range rankOver.AllRankOverFilter() {
		filter.Bounds = append(filter.Bounds, rankOverBound(bound))
	}
	return filter
}

func latestValue(ctx generated.ILatestFunctionContext) domain.FilterValue {
	raw := nodeText(ctx)
	value := domain.FilterValue{
		Kind:     domain.FilterValueLatest,
		Raw:      raw,
		Function: parseFunctionName(raw),
	}
	if ctx == nil {
		return value
	}
	for _, expression := range ctx.AllLatestExpression() {
		value.Arguments = append(value.Arguments, latestExpression(expression))
	}
	return value
}

func latestExpression(ctx generated.ILatestExpressionContext) domain.LatestExpression {
	expression := domain.LatestExpression{
		Raw:   nodeText(ctx),
		Field: terminalText(ctx.ID()),
	}
	if ctx == nil {
		return expression
	}
	if in := ctx.IN(); in != nil {
		expression.Operator = terminalText(in)
		expression.Value = timeIntervalValue(ctx.TimeIntervalArithmetic())
		return expression
	}

	expression.Operator = terminalText(ctx.COMPARISON_OPERATOR())
	switch {
	case ctx.PointInTimeArithmetic() != nil:
		expression.Value = pointInTimeValue(ctx.PointInTimeArithmetic())
	case ctx.TimeIntervalToPointInTime() != nil:
		expression.Value = intervalToPointInTimeValue(ctx.TimeIntervalToPointInTime())
	case ctx.SIGNED_INTEGER() != nil:
		expression.Value = domain.FilterValue{Kind: domain.FilterValueNumber, Raw: terminalText(ctx.SIGNED_INTEGER())}
	case ctx.FLOAT() != nil:
		expression.Value = domain.FilterValue{Kind: domain.FilterValueNumber, Raw: terminalText(ctx.FLOAT())}
	default:
		expression.Value = domain.FilterValue{Kind: domain.FilterValueGeneric, Raw: nodeText(ctx)}
	}
	return expression
}

func pointInTimeValue(ctx generated.IPointInTimeArithmeticContext) domain.FilterValue {
	raw := nodeText(ctx)
	value := domain.FilterValue{
		Kind:       domain.FilterValuePointInTime,
		Raw:        raw,
		Function:   parseFunctionName(raw),
		Arithmetic: arithmetic(ctx),
	}
	if ctx == nil {
		return value
	}
	value.TimeZone = timeZone(ctx.PointInTimeOrFunction())
	return value
}

func timeIntervalValue(ctx generated.ITimeIntervalArithmeticContext) domain.FilterValue {
	raw := nodeText(ctx)
	value := domain.FilterValue{
		Kind:       domain.FilterValueTimeInterval,
		Raw:        raw,
		Function:   parseFunctionName(raw),
		Arithmetic: arithmetic(ctx),
	}
	if ctx == nil {
		return value
	}
	if timeIntervalOrFunction := ctx.TimeIntervalOrFunction(); timeIntervalOrFunction != nil {
		value.TimeZone = timeZone(timeIntervalOrFunction)
		if interval := timeIntervalOrFunction.TimeInterval(); interval != nil {
			points := interval.AllPOINT_IN_TIME()
			if len(points) > 0 {
				value.Start = terminalText(points[0])
			}
			if len(points) > 1 {
				value.End = terminalText(points[1])
			}
		}
	}
	if gasIntervalOrFunction := ctx.GasIntervalOrFunction(); gasIntervalOrFunction != nil {
		value.TimeZone = timeZone(gasIntervalOrFunction)
	}
	return value
}

func intervalToPointInTimeValue(ctx generated.ITimeIntervalToPointInTimeContext) domain.FilterValue {
	raw := nodeText(ctx)
	return domain.FilterValue{
		Kind:     domain.FilterValueTimeIntervalPointTime,
		Raw:      raw,
		Function: parseFunctionName(raw),
	}
}

func rankOverBound(ctx generated.IRankOverFilterContext) domain.RankOverBound {
	bound := domain.RankOverBound{Raw: nodeText(ctx)}
	if ctx == nil {
		return bound
	}
	integers := ctx.AllSIGNED_INTEGER()
	if len(integers) > 0 {
		bound.Start = terminalText(integers[0])
	}
	if len(integers) > 1 {
		bound.End = terminalText(integers[1])
	}
	if ctx.OPEN_FILTER_INTERVAL_MARKER() != nil {
		bound.End = terminalText(ctx.OPEN_FILTER_INTERVAL_MARKER())
	}
	return bound
}

type arithmeticContext interface {
	GetArithmeticOperator() antlr.Token
	GetTimePeriod() antlr.Token
}

func arithmetic(ctx arithmeticContext) *domain.TimeArithmetic {
	if ctx == nil || ctx.GetArithmeticOperator() == nil {
		return nil
	}
	return &domain.TimeArithmetic{
		Operator: tokenText(ctx.GetArithmeticOperator()),
		Period:   tokenText(ctx.GetTimePeriod()),
	}
}

type textNode interface {
	GetText() string
}

func fieldName(id antlr.TerminalNode, surface generated.IKeySurfaceColumnContext) string {
	if id != nil {
		return terminalText(id)
	}
	return nodeText(surface)
}

func appendFilterNodes(nodes []domain.FilterNode, result interface{}) []domain.FilterNode {
	switch typed := result.(type) {
	case nil:
		return nodes
	case domain.FilterNode:
		return append(nodes, typed)
	case []domain.FilterNode:
		return append(nodes, typed...)
	default:
		return nodes
	}
}

func nodeText(node textNode) string {
	if node == nil {
		return ""
	}
	return node.GetText()
}

func terminalText(node antlr.TerminalNode) string {
	if node == nil {
		return ""
	}
	return node.GetText()
}

func tokenText(token antlr.Token) string {
	if token == nil {
		return ""
	}
	return token.GetText()
}

func terminalTexts(nodes []antlr.TerminalNode) []string {
	texts := make([]string, 0, len(nodes))
	for _, node := range nodes {
		texts = append(texts, terminalText(node))
	}
	return texts
}

type timeZoneContext interface {
	GetExpressionTimeZone() antlr.Token
}

func timeZone(ctx timeZoneContext) string {
	if ctx == nil {
		return ""
	}
	return tokenText(ctx.GetExpressionTimeZone())
}

func parseFunctionName(raw string) string {
	index := strings.Index(raw, "(")
	if index <= 0 {
		return ""
	}
	return raw[:index]
}
