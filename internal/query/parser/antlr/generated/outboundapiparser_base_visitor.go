// Code generated from OutboundAPIParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package generated // OutboundAPIParser
import "github.com/antlr4-go/antlr/v4"

type BaseOutboundAPIParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseOutboundAPIParserVisitor) VisitExpressionsSection(ctx *ExpressionsSectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitKeyFilterSection(ctx *KeyFilterSectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitIdPointInTimeArithmeticComparison(ctx *IdPointInTimeArithmeticComparisonContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitIdTimeIntervalIn(ctx *IdTimeIntervalInContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitIdNumericComparison(ctx *IdNumericComparisonContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitIdLatestGlobalComparison(ctx *IdLatestGlobalComparisonContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitIdTimeIntervalToPointInTimeComparison(ctx *IdTimeIntervalToPointInTimeComparisonContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitIdLatestComparison(ctx *IdLatestComparisonContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitTextComparison(ctx *TextComparisonContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitRankOver(ctx *RankOverContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitKeySurfaceColumn(ctx *KeySurfaceColumnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitTextColumn(ctx *TextColumnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitLatestGlobalFunction(ctx *LatestGlobalFunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitTimeInterval(ctx *TimeIntervalContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitTimeIntervalOrFunction(ctx *TimeIntervalOrFunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitGasIntervalOrFunction(ctx *GasIntervalOrFunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitPointInTimeOrFunction(ctx *PointInTimeOrFunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitPointInTimeArithmetic(ctx *PointInTimeArithmeticContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitTimeIntervalArithmetic(ctx *TimeIntervalArithmeticContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitTimeIntervalToPointInTime(ctx *TimeIntervalToPointInTimeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitRankOverFunction(ctx *RankOverFunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitRankOverFilter(ctx *RankOverFilterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitLatestFunction(ctx *LatestFunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitLatestExpression(ctx *LatestExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseOutboundAPIParserVisitor) VisitGenericValue(ctx *GenericValueContext) interface{} {
	return v.VisitChildren(ctx)
}
