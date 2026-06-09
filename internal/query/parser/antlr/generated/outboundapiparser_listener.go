// Code generated from OutboundAPIParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package generated // OutboundAPIParser
import "github.com/antlr4-go/antlr/v4"

// OutboundAPIParserListener is a complete listener for a parse tree produced by OutboundAPIParser.
type OutboundAPIParserListener interface {
	antlr.ParseTreeListener

	// EnterExpressionsSection is called when entering the expressionsSection production.
	EnterExpressionsSection(c *ExpressionsSectionContext)

	// EnterKeyFilterSection is called when entering the keyFilterSection production.
	EnterKeyFilterSection(c *KeyFilterSectionContext)

	// EnterIdPointInTimeArithmeticComparison is called when entering the IdPointInTimeArithmeticComparison production.
	EnterIdPointInTimeArithmeticComparison(c *IdPointInTimeArithmeticComparisonContext)

	// EnterIdTimeIntervalIn is called when entering the IdTimeIntervalIn production.
	EnterIdTimeIntervalIn(c *IdTimeIntervalInContext)

	// EnterIdNumericComparison is called when entering the IdNumericComparison production.
	EnterIdNumericComparison(c *IdNumericComparisonContext)

	// EnterIdLatestGlobalComparison is called when entering the IdLatestGlobalComparison production.
	EnterIdLatestGlobalComparison(c *IdLatestGlobalComparisonContext)

	// EnterIdTimeIntervalToPointInTimeComparison is called when entering the IdTimeIntervalToPointInTimeComparison production.
	EnterIdTimeIntervalToPointInTimeComparison(c *IdTimeIntervalToPointInTimeComparisonContext)

	// EnterIdLatestComparison is called when entering the IdLatestComparison production.
	EnterIdLatestComparison(c *IdLatestComparisonContext)

	// EnterTextComparison is called when entering the TextComparison production.
	EnterTextComparison(c *TextComparisonContext)

	// EnterRankOver is called when entering the RankOver production.
	EnterRankOver(c *RankOverContext)

	// EnterKeySurfaceColumn is called when entering the keySurfaceColumn production.
	EnterKeySurfaceColumn(c *KeySurfaceColumnContext)

	// EnterTextColumn is called when entering the textColumn production.
	EnterTextColumn(c *TextColumnContext)

	// EnterLatestGlobalFunction is called when entering the latestGlobalFunction production.
	EnterLatestGlobalFunction(c *LatestGlobalFunctionContext)

	// EnterTimeInterval is called when entering the timeInterval production.
	EnterTimeInterval(c *TimeIntervalContext)

	// EnterTimeIntervalOrFunction is called when entering the timeIntervalOrFunction production.
	EnterTimeIntervalOrFunction(c *TimeIntervalOrFunctionContext)

	// EnterGasIntervalOrFunction is called when entering the gasIntervalOrFunction production.
	EnterGasIntervalOrFunction(c *GasIntervalOrFunctionContext)

	// EnterPointInTimeOrFunction is called when entering the pointInTimeOrFunction production.
	EnterPointInTimeOrFunction(c *PointInTimeOrFunctionContext)

	// EnterPointInTimeArithmetic is called when entering the pointInTimeArithmetic production.
	EnterPointInTimeArithmetic(c *PointInTimeArithmeticContext)

	// EnterTimeIntervalArithmetic is called when entering the timeIntervalArithmetic production.
	EnterTimeIntervalArithmetic(c *TimeIntervalArithmeticContext)

	// EnterTimeIntervalToPointInTime is called when entering the timeIntervalToPointInTime production.
	EnterTimeIntervalToPointInTime(c *TimeIntervalToPointInTimeContext)

	// EnterRankOverFunction is called when entering the rankOverFunction production.
	EnterRankOverFunction(c *RankOverFunctionContext)

	// EnterRankOverFilter is called when entering the rankOverFilter production.
	EnterRankOverFilter(c *RankOverFilterContext)

	// EnterLatestFunction is called when entering the latestFunction production.
	EnterLatestFunction(c *LatestFunctionContext)

	// EnterLatestExpression is called when entering the latestExpression production.
	EnterLatestExpression(c *LatestExpressionContext)

	// EnterGenericValue is called when entering the genericValue production.
	EnterGenericValue(c *GenericValueContext)

	// ExitExpressionsSection is called when exiting the expressionsSection production.
	ExitExpressionsSection(c *ExpressionsSectionContext)

	// ExitKeyFilterSection is called when exiting the keyFilterSection production.
	ExitKeyFilterSection(c *KeyFilterSectionContext)

	// ExitIdPointInTimeArithmeticComparison is called when exiting the IdPointInTimeArithmeticComparison production.
	ExitIdPointInTimeArithmeticComparison(c *IdPointInTimeArithmeticComparisonContext)

	// ExitIdTimeIntervalIn is called when exiting the IdTimeIntervalIn production.
	ExitIdTimeIntervalIn(c *IdTimeIntervalInContext)

	// ExitIdNumericComparison is called when exiting the IdNumericComparison production.
	ExitIdNumericComparison(c *IdNumericComparisonContext)

	// ExitIdLatestGlobalComparison is called when exiting the IdLatestGlobalComparison production.
	ExitIdLatestGlobalComparison(c *IdLatestGlobalComparisonContext)

	// ExitIdTimeIntervalToPointInTimeComparison is called when exiting the IdTimeIntervalToPointInTimeComparison production.
	ExitIdTimeIntervalToPointInTimeComparison(c *IdTimeIntervalToPointInTimeComparisonContext)

	// ExitIdLatestComparison is called when exiting the IdLatestComparison production.
	ExitIdLatestComparison(c *IdLatestComparisonContext)

	// ExitTextComparison is called when exiting the TextComparison production.
	ExitTextComparison(c *TextComparisonContext)

	// ExitRankOver is called when exiting the RankOver production.
	ExitRankOver(c *RankOverContext)

	// ExitKeySurfaceColumn is called when exiting the keySurfaceColumn production.
	ExitKeySurfaceColumn(c *KeySurfaceColumnContext)

	// ExitTextColumn is called when exiting the textColumn production.
	ExitTextColumn(c *TextColumnContext)

	// ExitLatestGlobalFunction is called when exiting the latestGlobalFunction production.
	ExitLatestGlobalFunction(c *LatestGlobalFunctionContext)

	// ExitTimeInterval is called when exiting the timeInterval production.
	ExitTimeInterval(c *TimeIntervalContext)

	// ExitTimeIntervalOrFunction is called when exiting the timeIntervalOrFunction production.
	ExitTimeIntervalOrFunction(c *TimeIntervalOrFunctionContext)

	// ExitGasIntervalOrFunction is called when exiting the gasIntervalOrFunction production.
	ExitGasIntervalOrFunction(c *GasIntervalOrFunctionContext)

	// ExitPointInTimeOrFunction is called when exiting the pointInTimeOrFunction production.
	ExitPointInTimeOrFunction(c *PointInTimeOrFunctionContext)

	// ExitPointInTimeArithmetic is called when exiting the pointInTimeArithmetic production.
	ExitPointInTimeArithmetic(c *PointInTimeArithmeticContext)

	// ExitTimeIntervalArithmetic is called when exiting the timeIntervalArithmetic production.
	ExitTimeIntervalArithmetic(c *TimeIntervalArithmeticContext)

	// ExitTimeIntervalToPointInTime is called when exiting the timeIntervalToPointInTime production.
	ExitTimeIntervalToPointInTime(c *TimeIntervalToPointInTimeContext)

	// ExitRankOverFunction is called when exiting the rankOverFunction production.
	ExitRankOverFunction(c *RankOverFunctionContext)

	// ExitRankOverFilter is called when exiting the rankOverFilter production.
	ExitRankOverFilter(c *RankOverFilterContext)

	// ExitLatestFunction is called when exiting the latestFunction production.
	ExitLatestFunction(c *LatestFunctionContext)

	// ExitLatestExpression is called when exiting the latestExpression production.
	ExitLatestExpression(c *LatestExpressionContext)

	// ExitGenericValue is called when exiting the genericValue production.
	ExitGenericValue(c *GenericValueContext)
}
