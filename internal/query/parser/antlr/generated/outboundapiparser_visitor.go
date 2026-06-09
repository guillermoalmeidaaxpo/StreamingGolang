// Code generated from OutboundAPIParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package generated // OutboundAPIParser
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by OutboundAPIParser.
type OutboundAPIParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by OutboundAPIParser#expressionsSection.
	VisitExpressionsSection(ctx *ExpressionsSectionContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#keyFilterSection.
	VisitKeyFilterSection(ctx *KeyFilterSectionContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#IdPointInTimeArithmeticComparison.
	VisitIdPointInTimeArithmeticComparison(ctx *IdPointInTimeArithmeticComparisonContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#IdTimeIntervalIn.
	VisitIdTimeIntervalIn(ctx *IdTimeIntervalInContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#IdNumericComparison.
	VisitIdNumericComparison(ctx *IdNumericComparisonContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#IdLatestGlobalComparison.
	VisitIdLatestGlobalComparison(ctx *IdLatestGlobalComparisonContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#IdTimeIntervalToPointInTimeComparison.
	VisitIdTimeIntervalToPointInTimeComparison(ctx *IdTimeIntervalToPointInTimeComparisonContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#IdLatestComparison.
	VisitIdLatestComparison(ctx *IdLatestComparisonContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#TextComparison.
	VisitTextComparison(ctx *TextComparisonContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#RankOver.
	VisitRankOver(ctx *RankOverContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#keySurfaceColumn.
	VisitKeySurfaceColumn(ctx *KeySurfaceColumnContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#textColumn.
	VisitTextColumn(ctx *TextColumnContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#latestGlobalFunction.
	VisitLatestGlobalFunction(ctx *LatestGlobalFunctionContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#timeInterval.
	VisitTimeInterval(ctx *TimeIntervalContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#timeIntervalOrFunction.
	VisitTimeIntervalOrFunction(ctx *TimeIntervalOrFunctionContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#gasIntervalOrFunction.
	VisitGasIntervalOrFunction(ctx *GasIntervalOrFunctionContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#pointInTimeOrFunction.
	VisitPointInTimeOrFunction(ctx *PointInTimeOrFunctionContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#pointInTimeArithmetic.
	VisitPointInTimeArithmetic(ctx *PointInTimeArithmeticContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#timeIntervalArithmetic.
	VisitTimeIntervalArithmetic(ctx *TimeIntervalArithmeticContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#timeIntervalToPointInTime.
	VisitTimeIntervalToPointInTime(ctx *TimeIntervalToPointInTimeContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#rankOverFunction.
	VisitRankOverFunction(ctx *RankOverFunctionContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#rankOverFilter.
	VisitRankOverFilter(ctx *RankOverFilterContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#latestFunction.
	VisitLatestFunction(ctx *LatestFunctionContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#latestExpression.
	VisitLatestExpression(ctx *LatestExpressionContext) interface{}

	// Visit a parse tree produced by OutboundAPIParser#genericValue.
	VisitGenericValue(ctx *GenericValueContext) interface{}
}
