parser grammar OutboundAPIParser;

options
{
    tokenVocab = OutboundAPILexer;
}

expressionsSection
    : keyFilterSection+
    ;

keyFilterSection
    :  keyComparison (SEMICOLON keyComparison )* SEMICOLON? EOF
    ;

keyComparison
    : (ID | keySurfaceColumn) COMPARISON_OPERATOR pointInTimeArithmetic             #IdPointInTimeArithmeticComparison
    | (ID | keySurfaceColumn) IN timeIntervalArithmetic                             #IdTimeIntervalIn
    | ID COMPARISON_OPERATOR number=(SIGNED_INTEGER | FLOAT)                        #IdNumericComparison
    | ID COMPARISON_OPERATOR latestGlobalFunction                                   #IdLatestGlobalComparison
    | (ID | keySurfaceColumn) COMPARISON_OPERATOR timeIntervalToPointInTime         #IdTimeIntervalToPointInTimeComparison
    | ID COMPARISON_OPERATOR latestFunction                                         #IdLatestComparison
    | textColumn COMPARISON_OPERATOR genericValue                                   #TextComparison
    | rankOverFunction                                                              #RankOver
    ;

keySurfaceColumn
    : CI_VALIDITY_PERIOD_START
    | CI_VALIDITY_PERIOD_END
    | CI_INSTANCE_CODE
    | CI_ISP
    ;

textColumn
    : CI_BIDID
    | CI_DIRECTION
    | CI_STATUS
    | CI_OPTION_EXPIRY
    ;

latestGlobalFunction
    : LATEST_GLOBAL LB RB
    ;

timeInterval
    : POINT_IN_TIME COMMA POINT_IN_TIME
    ;

timeIntervalOrFunction
    : TIME_INTERVAL_EXPLICIT_FUNCTION_NAME LB timeInterval RB
    | TIME_INTERVAL_FUNCTION_NAME LB pointInTimeArithmetic (COMMA expressionTimeZone = TIME_ZONE_IANA)? RB
    ;

gasIntervalOrFunction
    : TIME_INTERVAL_GAS_FUNCTION_NAME LB pointInTimeArithmetic (COMMA expressionTimeZone = TIME_ZONE_IANA)? RB
    ;

pointInTimeOrFunction
    : POINT_IN_TIME
    | POINT_IN_TIME_FUNCTION_NAME LB RB
    | POINT_IN_TIME_UTC_FUNCTION_NAME LB (POINT_IN_TIME | POINT_IN_TIME_FUNCTION_NAME LB RB) COMMA expressionTimeZone = TIME_ZONE_IANA RB
    ;

pointInTimeArithmetic
    : (pointInTimeOrFunction) (ArithmeticOperator = (ADD | SUB) TimePeriod=TIME_PERIOD)?
    ;

timeIntervalArithmetic
    : (timeIntervalOrFunction) (ArithmeticOperator = (ADD | SUB) TimePeriod=TIME_PERIOD)?
    | gasIntervalOrFunction
    ;

timeIntervalToPointInTime
    : TIME_INTERVAL_TO_POINT_IN_TIME_FUNCTION LB timeIntervalArithmetic RB
    ;

rankOverFunction
    : RANK_OVER_FUNCTION_NAME LB LSB ID (COMMA ID)* RSB COMMA LSB ID SORT_ORDER (COMMA ID SORT_ORDER)* RSB (COMMA rankOverFilter)* RB
    ;

rankOverFilter
    : SIGNED_INTEGER
    | LSB SIGNED_INTEGER COMMA SIGNED_INTEGER RSB
    | LSB SIGNED_INTEGER COMMA OPEN_FILTER_INTERVAL_MARKER RSB
    ;

latestFunction
    : LATEST_FUNCTION_NAME LB latestExpression (COMMA latestExpression)* RB
    ;

latestExpression
    : ID COMPARISON_OPERATOR pointInTimeArithmetic
    | ID IN timeIntervalArithmetic
    | ID COMPARISON_OPERATOR timeIntervalToPointInTime
    | ID COMPARISON_OPERATOR (SIGNED_INTEGER | FLOAT)
    ;

genericValue
    : IN
    | SORT_ORDER
    | OPEN_FILTER_INTERVAL_MARKER
    | DATE
    | TIME
    | POINT_IN_TIME
    | TIME_PERIOD
    | TIME_ZONE_IANA
    | SIGNED_INTEGER
    | FLOAT
    | ID
    | WORD
    ;
