// Code generated from OutboundAPIParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package generated // OutboundAPIParser
import "github.com/antlr4-go/antlr/v4"

// BaseOutboundAPIParserListener is a complete listener for a parse tree produced by OutboundAPIParser.
type BaseOutboundAPIParserListener struct{}

var _ OutboundAPIParserListener = &BaseOutboundAPIParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseOutboundAPIParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseOutboundAPIParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseOutboundAPIParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseOutboundAPIParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterExpressionsSection is called when production expressionsSection is entered.
func (s *BaseOutboundAPIParserListener) EnterExpressionsSection(ctx *ExpressionsSectionContext) {}

// ExitExpressionsSection is called when production expressionsSection is exited.
func (s *BaseOutboundAPIParserListener) ExitExpressionsSection(ctx *ExpressionsSectionContext) {}

// EnterKeyFilterSection is called when production keyFilterSection is entered.
func (s *BaseOutboundAPIParserListener) EnterKeyFilterSection(ctx *KeyFilterSectionContext) {}

// ExitKeyFilterSection is called when production keyFilterSection is exited.
func (s *BaseOutboundAPIParserListener) ExitKeyFilterSection(ctx *KeyFilterSectionContext) {}

// EnterIdPointInTimeArithmeticComparison is called when production IdPointInTimeArithmeticComparison is entered.
func (s *BaseOutboundAPIParserListener) EnterIdPointInTimeArithmeticComparison(ctx *IdPointInTimeArithmeticComparisonContext) {
}

// ExitIdPointInTimeArithmeticComparison is called when production IdPointInTimeArithmeticComparison is exited.
func (s *BaseOutboundAPIParserListener) ExitIdPointInTimeArithmeticComparison(ctx *IdPointInTimeArithmeticComparisonContext) {
}

// EnterIdTimeIntervalIn is called when production IdTimeIntervalIn is entered.
func (s *BaseOutboundAPIParserListener) EnterIdTimeIntervalIn(ctx *IdTimeIntervalInContext) {}

// ExitIdTimeIntervalIn is called when production IdTimeIntervalIn is exited.
func (s *BaseOutboundAPIParserListener) ExitIdTimeIntervalIn(ctx *IdTimeIntervalInContext) {}

// EnterIdNumericComparison is called when production IdNumericComparison is entered.
func (s *BaseOutboundAPIParserListener) EnterIdNumericComparison(ctx *IdNumericComparisonContext) {}

// ExitIdNumericComparison is called when production IdNumericComparison is exited.
func (s *BaseOutboundAPIParserListener) ExitIdNumericComparison(ctx *IdNumericComparisonContext) {}

// EnterIdLatestGlobalComparison is called when production IdLatestGlobalComparison is entered.
func (s *BaseOutboundAPIParserListener) EnterIdLatestGlobalComparison(ctx *IdLatestGlobalComparisonContext) {
}

// ExitIdLatestGlobalComparison is called when production IdLatestGlobalComparison is exited.
func (s *BaseOutboundAPIParserListener) ExitIdLatestGlobalComparison(ctx *IdLatestGlobalComparisonContext) {
}

// EnterIdTimeIntervalToPointInTimeComparison is called when production IdTimeIntervalToPointInTimeComparison is entered.
func (s *BaseOutboundAPIParserListener) EnterIdTimeIntervalToPointInTimeComparison(ctx *IdTimeIntervalToPointInTimeComparisonContext) {
}

// ExitIdTimeIntervalToPointInTimeComparison is called when production IdTimeIntervalToPointInTimeComparison is exited.
func (s *BaseOutboundAPIParserListener) ExitIdTimeIntervalToPointInTimeComparison(ctx *IdTimeIntervalToPointInTimeComparisonContext) {
}

// EnterIdLatestComparison is called when production IdLatestComparison is entered.
func (s *BaseOutboundAPIParserListener) EnterIdLatestComparison(ctx *IdLatestComparisonContext) {}

// ExitIdLatestComparison is called when production IdLatestComparison is exited.
func (s *BaseOutboundAPIParserListener) ExitIdLatestComparison(ctx *IdLatestComparisonContext) {}

// EnterTextComparison is called when production TextComparison is entered.
func (s *BaseOutboundAPIParserListener) EnterTextComparison(ctx *TextComparisonContext) {}

// ExitTextComparison is called when production TextComparison is exited.
func (s *BaseOutboundAPIParserListener) ExitTextComparison(ctx *TextComparisonContext) {}

// EnterRankOver is called when production RankOver is entered.
func (s *BaseOutboundAPIParserListener) EnterRankOver(ctx *RankOverContext) {}

// ExitRankOver is called when production RankOver is exited.
func (s *BaseOutboundAPIParserListener) ExitRankOver(ctx *RankOverContext) {}

// EnterKeySurfaceColumn is called when production keySurfaceColumn is entered.
func (s *BaseOutboundAPIParserListener) EnterKeySurfaceColumn(ctx *KeySurfaceColumnContext) {}

// ExitKeySurfaceColumn is called when production keySurfaceColumn is exited.
func (s *BaseOutboundAPIParserListener) ExitKeySurfaceColumn(ctx *KeySurfaceColumnContext) {}

// EnterTextColumn is called when production textColumn is entered.
func (s *BaseOutboundAPIParserListener) EnterTextColumn(ctx *TextColumnContext) {}

// ExitTextColumn is called when production textColumn is exited.
func (s *BaseOutboundAPIParserListener) ExitTextColumn(ctx *TextColumnContext) {}

// EnterLatestGlobalFunction is called when production latestGlobalFunction is entered.
func (s *BaseOutboundAPIParserListener) EnterLatestGlobalFunction(ctx *LatestGlobalFunctionContext) {}

// ExitLatestGlobalFunction is called when production latestGlobalFunction is exited.
func (s *BaseOutboundAPIParserListener) ExitLatestGlobalFunction(ctx *LatestGlobalFunctionContext) {}

// EnterTimeInterval is called when production timeInterval is entered.
func (s *BaseOutboundAPIParserListener) EnterTimeInterval(ctx *TimeIntervalContext) {}

// ExitTimeInterval is called when production timeInterval is exited.
func (s *BaseOutboundAPIParserListener) ExitTimeInterval(ctx *TimeIntervalContext) {}

// EnterTimeIntervalOrFunction is called when production timeIntervalOrFunction is entered.
func (s *BaseOutboundAPIParserListener) EnterTimeIntervalOrFunction(ctx *TimeIntervalOrFunctionContext) {
}

// ExitTimeIntervalOrFunction is called when production timeIntervalOrFunction is exited.
func (s *BaseOutboundAPIParserListener) ExitTimeIntervalOrFunction(ctx *TimeIntervalOrFunctionContext) {
}

// EnterGasIntervalOrFunction is called when production gasIntervalOrFunction is entered.
func (s *BaseOutboundAPIParserListener) EnterGasIntervalOrFunction(ctx *GasIntervalOrFunctionContext) {
}

// ExitGasIntervalOrFunction is called when production gasIntervalOrFunction is exited.
func (s *BaseOutboundAPIParserListener) ExitGasIntervalOrFunction(ctx *GasIntervalOrFunctionContext) {
}

// EnterPointInTimeOrFunction is called when production pointInTimeOrFunction is entered.
func (s *BaseOutboundAPIParserListener) EnterPointInTimeOrFunction(ctx *PointInTimeOrFunctionContext) {
}

// ExitPointInTimeOrFunction is called when production pointInTimeOrFunction is exited.
func (s *BaseOutboundAPIParserListener) ExitPointInTimeOrFunction(ctx *PointInTimeOrFunctionContext) {
}

// EnterPointInTimeArithmetic is called when production pointInTimeArithmetic is entered.
func (s *BaseOutboundAPIParserListener) EnterPointInTimeArithmetic(ctx *PointInTimeArithmeticContext) {
}

// ExitPointInTimeArithmetic is called when production pointInTimeArithmetic is exited.
func (s *BaseOutboundAPIParserListener) ExitPointInTimeArithmetic(ctx *PointInTimeArithmeticContext) {
}

// EnterTimeIntervalArithmetic is called when production timeIntervalArithmetic is entered.
func (s *BaseOutboundAPIParserListener) EnterTimeIntervalArithmetic(ctx *TimeIntervalArithmeticContext) {
}

// ExitTimeIntervalArithmetic is called when production timeIntervalArithmetic is exited.
func (s *BaseOutboundAPIParserListener) ExitTimeIntervalArithmetic(ctx *TimeIntervalArithmeticContext) {
}

// EnterTimeIntervalToPointInTime is called when production timeIntervalToPointInTime is entered.
func (s *BaseOutboundAPIParserListener) EnterTimeIntervalToPointInTime(ctx *TimeIntervalToPointInTimeContext) {
}

// ExitTimeIntervalToPointInTime is called when production timeIntervalToPointInTime is exited.
func (s *BaseOutboundAPIParserListener) ExitTimeIntervalToPointInTime(ctx *TimeIntervalToPointInTimeContext) {
}

// EnterRankOverFunction is called when production rankOverFunction is entered.
func (s *BaseOutboundAPIParserListener) EnterRankOverFunction(ctx *RankOverFunctionContext) {}

// ExitRankOverFunction is called when production rankOverFunction is exited.
func (s *BaseOutboundAPIParserListener) ExitRankOverFunction(ctx *RankOverFunctionContext) {}

// EnterRankOverFilter is called when production rankOverFilter is entered.
func (s *BaseOutboundAPIParserListener) EnterRankOverFilter(ctx *RankOverFilterContext) {}

// ExitRankOverFilter is called when production rankOverFilter is exited.
func (s *BaseOutboundAPIParserListener) ExitRankOverFilter(ctx *RankOverFilterContext) {}

// EnterLatestFunction is called when production latestFunction is entered.
func (s *BaseOutboundAPIParserListener) EnterLatestFunction(ctx *LatestFunctionContext) {}

// ExitLatestFunction is called when production latestFunction is exited.
func (s *BaseOutboundAPIParserListener) ExitLatestFunction(ctx *LatestFunctionContext) {}

// EnterLatestExpression is called when production latestExpression is entered.
func (s *BaseOutboundAPIParserListener) EnterLatestExpression(ctx *LatestExpressionContext) {}

// ExitLatestExpression is called when production latestExpression is exited.
func (s *BaseOutboundAPIParserListener) ExitLatestExpression(ctx *LatestExpressionContext) {}

// EnterGenericValue is called when production genericValue is entered.
func (s *BaseOutboundAPIParserListener) EnterGenericValue(ctx *GenericValueContext) {}

// ExitGenericValue is called when production genericValue is exited.
func (s *BaseOutboundAPIParserListener) ExitGenericValue(ctx *GenericValueContext) {}
