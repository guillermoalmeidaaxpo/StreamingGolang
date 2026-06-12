package antlrparser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antlr4-go/antlr/v4"

	"streaming-golang/internal/domain"
	"streaming-golang/internal/domain/timeexpr"
	"streaming-golang/internal/query/parser/antlr/generated"
)

const (
	onlyReferenceTimeColumnAllowedMessage       = "Only 'ReferenceTime' can be used as the target of the 'latest' and 'latestGlobal' functions."
	notAllowedLatestOperatorMessage             = "Only equality comparisons ('=') are allowed with 'latest' and 'latestGlobal' functions."
	latestExpressionParameterCountMessage       = "'latest' function only accepts a single expression parameter."
	onlyReferenceTimeAllowedInsideLatestMessage = "Only 'ReferenceTime' can be used in the expression inside of the 'latest' function."
)

var isoPeriodPattern = regexp.MustCompile(`^P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)W)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?)?$`)

type astVisitor struct {
	*generated.BaseOutboundAPIParserVisitor
	timeZone string
	errors   []string
}

func newASTVisitor(timeZone string) *astVisitor {
	return &astVisitor{
		BaseOutboundAPIParserVisitor: &generated.BaseOutboundAPIParserVisitor{
			BaseParseTreeVisitor: &antlr.BaseParseTreeVisitor{},
		},
		timeZone: timeZone,
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
		Value:    v.pointInTimeValue(ctx.PointInTimeArithmetic()),
	}
}

func (v *astVisitor) VisitIdTimeIntervalIn(ctx *generated.IdTimeIntervalInContext) interface{} {
	field := fieldName(ctx.ID(), ctx.KeySurfaceColumn())
	value := v.timeIntervalValue(ctx.TimeIntervalArithmetic())
	if value.Start != "" && value.End != "" {
		return []domain.FilterNode{
			domain.ComparisonFilter{
				Raw:      ctx.GetText(),
				Field:    field,
				Operator: ">=",
				Value:    domain.FilterValue{Kind: domain.FilterValuePointInTime, Raw: value.Start, TimeZone: value.TimeZone},
			},
			domain.ComparisonFilter{
				Raw:      ctx.GetText(),
				Field:    field,
				Operator: "<",
				Value:    domain.FilterValue{Kind: domain.FilterValuePointInTime, Raw: value.End, TimeZone: value.TimeZone},
			},
		}
	}
	return domain.ComparisonFilter{Raw: ctx.GetText(), Field: field, Operator: terminalText(ctx.IN()), Value: value}
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
	field := terminalText(ctx.ID())
	operator := terminalText(ctx.COMPARISON_OPERATOR())
	if !strings.EqualFold(field, "ReferenceTime") {
		v.errors = append(v.errors, onlyReferenceTimeColumnAllowedMessage)
	}
	if !strings.EqualFold(operator, "=") {
		v.errors = append(v.errors, notAllowedLatestOperatorMessage)
	}
	raw := nodeText(ctx.LatestGlobalFunction())
	return domain.ComparisonFilter{
		Raw:      ctx.GetText(),
		Field:    field,
		Operator: operator,
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
		Value:    v.intervalToPointInTimeValue(ctx.TimeIntervalToPointInTime()),
	}
}

func (v *astVisitor) VisitIdLatestComparison(ctx *generated.IdLatestComparisonContext) interface{} {
	field := terminalText(ctx.ID())
	operator := terminalText(ctx.COMPARISON_OPERATOR())
	if !strings.EqualFold(field, "ReferenceTime") {
		v.errors = append(v.errors, onlyReferenceTimeColumnAllowedMessage)
	}
	if !strings.EqualFold(operator, "=") {
		v.errors = append(v.errors, notAllowedLatestOperatorMessage)
	}
	if latest := ctx.LatestFunction(); latest != nil {
		expressions := latest.AllLatestExpression()
		if len(expressions) != 1 {
			v.errors = append(v.errors, latestExpressionParameterCountMessage)
		} else if !strings.EqualFold(terminalText(expressions[0].ID()), "ReferenceTime") {
			v.errors = append(v.errors, onlyReferenceTimeAllowedInsideLatestMessage)
		}
	}
	return domain.ComparisonFilter{
		Raw:      ctx.GetText(),
		Field:    field,
		Operator: operator,
		Value:    v.latestValue(ctx.LatestFunction()),
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

func (v *astVisitor) latestValue(ctx generated.ILatestFunctionContext) domain.FilterValue {
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
		value.Arguments = append(value.Arguments, v.latestExpression(expression))
	}
	return value
}

func (v *astVisitor) latestExpression(ctx generated.ILatestExpressionContext) domain.LatestExpression {
	expression := domain.LatestExpression{
		Raw:   nodeText(ctx),
		Field: terminalText(ctx.ID()),
	}
	if ctx == nil {
		return expression
	}
	if in := ctx.IN(); in != nil {
		expression.Operator = terminalText(in)
		expression.Value = v.timeIntervalValue(ctx.TimeIntervalArithmetic())
		return expression
	}

	expression.Operator = terminalText(ctx.COMPARISON_OPERATOR())
	switch {
	case ctx.PointInTimeArithmetic() != nil:
		expression.Value = v.pointInTimeValue(ctx.PointInTimeArithmetic())
	case ctx.TimeIntervalToPointInTime() != nil:
		expression.Value = v.intervalToPointInTimeValue(ctx.TimeIntervalToPointInTime())
	case ctx.SIGNED_INTEGER() != nil:
		expression.Value = domain.FilterValue{Kind: domain.FilterValueNumber, Raw: terminalText(ctx.SIGNED_INTEGER())}
	case ctx.FLOAT() != nil:
		expression.Value = domain.FilterValue{Kind: domain.FilterValueNumber, Raw: terminalText(ctx.FLOAT())}
	default:
		expression.Value = domain.FilterValue{Kind: domain.FilterValueGeneric, Raw: nodeText(ctx)}
	}
	return expression
}

func (v *astVisitor) pointInTimeValue(ctx generated.IPointInTimeArithmeticContext) domain.FilterValue {
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
	pointInTime, ok := v.parsePointInTime(ctx.PointInTimeOrFunction())
	if ok {
		if value.Arithmetic != nil {
			var err error
			pointInTime, err = applyPeriod(pointInTime, value.Arithmetic.Operator, value.Arithmetic.Period)
			if err != nil {
				v.errors = append(v.errors, err.Error())
			}
		}
		value.Raw = timeexpr.FormatUTC(pointInTime)
	}
	value.TimeZone = timeZone(ctx.PointInTimeOrFunction())
	return value
}

func (v *astVisitor) timeIntervalValue(ctx generated.ITimeIntervalArithmeticContext) domain.FilterValue {
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
				value.Start = v.normalizedPointInTimeToken(terminalText(points[0]), v.effectiveTimeZone(value.TimeZone))
			}
			if len(points) > 1 {
				value.End = v.normalizedPointInTimeToken(terminalText(points[1]), v.effectiveTimeZone(value.TimeZone))
			}
		} else if start, end, ok := v.timeIntervalFunctionBounds(timeIntervalOrFunction.GetText(), value.TimeZone); ok {
			value.Start = timeexpr.FormatUTC(start)
			value.End = timeexpr.FormatUTC(end)
		}
	}
	if gasIntervalOrFunction := ctx.GasIntervalOrFunction(); gasIntervalOrFunction != nil {
		value.TimeZone = timeZone(gasIntervalOrFunction)
		if start, end, ok := v.timeIntervalFunctionBounds(gasIntervalOrFunction.GetText(), value.TimeZone); ok {
			value.Start = timeexpr.FormatUTC(start)
			value.End = timeexpr.FormatUTC(end)
		}
	}
	if value.Start != "" && value.End != "" && value.Arithmetic != nil {
		start, startOK := v.parsePointInTimeRaw(value.Start)
		end, endOK := v.parsePointInTimeRaw(value.End)
		if startOK && endOK {
			shiftedStart, err := applyPeriod(start, value.Arithmetic.Operator, value.Arithmetic.Period)
			if err != nil {
				v.errors = append(v.errors, err.Error())
			}
			shiftedEnd, err := applyPeriod(end, value.Arithmetic.Operator, value.Arithmetic.Period)
			if err != nil {
				v.errors = append(v.errors, err.Error())
			}
			value.Start = timeexpr.FormatUTC(shiftedStart)
			value.End = timeexpr.FormatUTC(shiftedEnd)
		}
	}
	return value
}

func (v *astVisitor) intervalToPointInTimeValue(ctx generated.ITimeIntervalToPointInTimeContext) domain.FilterValue {
	raw := nodeText(ctx)
	value := domain.FilterValue{
		Kind:     domain.FilterValueTimeIntervalPointTime,
		Raw:      raw,
		Function: parseFunctionName(raw),
	}
	if ctx == nil {
		return value
	}
	interval := v.timeIntervalValue(ctx.TimeIntervalArithmetic())
	if interval.Start == "" || interval.End == "" {
		return value
	}
	switch strings.ToLower(value.Function) {
	case "begin":
		value.Raw = interval.Start
		value.Kind = domain.FilterValuePointInTime
	case "end":
		value.Raw = interval.End
		value.Kind = domain.FilterValuePointInTime
	}
	return domain.FilterValue{
		Kind:     value.Kind,
		Raw:      value.Raw,
		Function: value.Function,
	}
}

func (v *astVisitor) parsePointInTime(ctx generated.IPointInTimeOrFunctionContext) (time.Time, bool) {
	if ctx == nil {
		return time.Time{}, false
	}
	if point := ctx.POINT_IN_TIME(); point != nil {
		return v.parsePointInTimeToken(terminalText(point), v.effectiveTimeZone(""))
	}
	if ctx.POINT_IN_TIME_FUNCTION_NAME() != nil {
		return time.Now().UTC(), true
	}
	if ctx.POINT_IN_TIME_UTC_FUNCTION_NAME() != nil {
		zone := terminalText(ctx.TIME_ZONE_IANA())
		if point := ctx.POINT_IN_TIME(); point != nil {
			return v.parsePointInTimeToken(terminalText(point), zone)
		}
		return time.Now().UTC(), true
	}
	return time.Time{}, false
}

func (v *astVisitor) normalizedPointInTimeToken(raw string, timeZone string) string {
	point, ok := v.parsePointInTimeToken(raw, timeZone)
	if !ok {
		return raw
	}
	return timeexpr.FormatUTC(point)
}

func (v *astVisitor) parsePointInTimeToken(raw string, timeZone string) (time.Time, bool) {
	loc, err := timeexpr.LoadLocation(v.effectiveTimeZone(timeZone))
	if err != nil {
		v.errors = append(v.errors, fmt.Sprintf("invalid timezone %q", v.effectiveTimeZone(timeZone)))
		return time.Time{}, false
	}
	point, err := timeexpr.ParsePointInTimeToken(raw, loc)
	if err != nil {
		v.errors = append(v.errors, err.Error())
		return time.Time{}, false
	}
	return point, true
}

func (v *astVisitor) parsePointInTimeRaw(raw string) (time.Time, bool) {
	point, err := timeexpr.ParsePointInTime(raw, time.UTC)
	if err != nil {
		v.errors = append(v.errors, err.Error())
		return time.Time{}, false
	}
	return point, true
}

func (v *astVisitor) timeIntervalFunctionBounds(raw string, expressionTimeZone string) (time.Time, time.Time, bool) {
	name, args, ok := functionCall(raw)
	if !ok {
		return time.Time{}, time.Time{}, false
	}
	parts := splitArguments(args)
	if len(parts) == 0 {
		return time.Time{}, time.Time{}, false
	}

	loc, err := timeexpr.LoadLocation(v.effectiveTimeZone(expressionTimeZone))
	if err != nil {
		v.errors = append(v.errors, fmt.Sprintf("invalid timezone %q", v.effectiveTimeZone(expressionTimeZone)))
		return time.Time{}, time.Time{}, false
	}

	start, err := parsePointInTimeWithArithmetic(parts[0], loc)
	if err != nil {
		v.errors = append(v.errors, err.Error())
		return time.Time{}, time.Time{}, false
	}

	switch strings.ToLower(name) {
	case "tiday", "gasdayeurope":
		return start, start.AddDate(0, 0, 1), true
	case "tiweek", "gasweekeurope":
		return start, start.AddDate(0, 0, 7), true
	case "timonth", "gasmontheurope":
		return start, start.AddDate(0, 1, 0), true
	case "tiquarter", "gasquartereurope":
		return start, start.AddDate(0, 3, 0), true
	case "tiyear", "gasyeareurope":
		return start, start.AddDate(1, 0, 0), true
	case "gassummereurope":
		year := start.In(loc).Year()
		return time.Date(year, time.April, 1, 6, 0, 0, 0, loc).UTC(), time.Date(year, time.October, 1, 6, 0, 0, 0, loc).UTC(), true
	case "gaswintereurope":
		localStart := start.In(loc)
		year := localStart.Year()
		if localStart.Month() < time.October {
			year--
		}
		return time.Date(year, time.October, 1, 6, 0, 0, 0, loc).UTC(), time.Date(year+1, time.April, 1, 6, 0, 0, 0, loc).UTC(), true
	default:
		return time.Time{}, time.Time{}, false
	}
}

func (v *astVisitor) effectiveTimeZone(timeZone string) string {
	if strings.TrimSpace(timeZone) != "" {
		return timeZone
	}
	return v.timeZone
}

func (v *astVisitor) hasErrors() bool {
	return len(v.errors) > 0
}

func (v *astVisitor) message() string {
	return strings.Join(v.errors, "\n")
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

func parsePointInTimeWithArithmetic(raw string, loc *time.Location) (time.Time, error) {
	base, operator, period := splitPointTimeArithmetic(raw)
	var point time.Time
	var err error
	if strings.EqualFold(strings.TrimSpace(base), "now()") {
		point = time.Now().In(loc).UTC()
	} else {
		point, err = timeexpr.ParsePointInTime(base, loc)
		if err != nil {
			return time.Time{}, err
		}
	}
	if operator == "" {
		return point, nil
	}
	return applyPeriod(point, operator, period)
}

func splitPointTimeArithmetic(raw string) (base, operator, period string) {
	for _, marker := range []string{"+P", "-P"} {
		if index := strings.Index(raw, marker); index > 0 {
			return raw[:index], raw[index : index+1], raw[index+1:]
		}
	}
	return raw, "", ""
}

func applyPeriod(value time.Time, operator, rawPeriod string) (time.Time, error) {
	parts := isoPeriodPattern.FindStringSubmatch(rawPeriod)
	if parts == nil {
		return time.Time{}, fmt.Errorf("invalid ISO period %q", rawPeriod)
	}

	sign := 1
	if operator == "-" {
		sign = -1
	}

	years := atoi(parts[1]) * sign
	months := atoi(parts[2]) * sign
	weeks := atoi(parts[3])
	days := (atoi(parts[4]) + weeks*7) * sign
	hours := atoi(parts[5]) * sign
	minutes := atoi(parts[6]) * sign
	seconds := atoi(parts[7]) * sign

	value = value.AddDate(years, months, days)
	return value.Add(time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second), nil
}

func atoi(raw string) int {
	if raw == "" {
		return 0
	}
	value, _ := strconv.Atoi(raw)
	return value
}

func functionCall(raw string) (name, args string, ok bool) {
	raw = strings.TrimSpace(raw)
	open := strings.Index(raw, "(")
	if open <= 0 || !strings.HasSuffix(raw, ")") {
		return "", "", false
	}
	return raw[:open], raw[open+1 : len(raw)-1], true
}

func splitArguments(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
