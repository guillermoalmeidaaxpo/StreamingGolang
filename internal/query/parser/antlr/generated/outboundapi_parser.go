// Code generated from OutboundAPIParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package generated // OutboundAPIParser
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type OutboundAPIParser struct {
	*antlr.BaseParser
}

var OutboundAPIParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func outboundapiparserParserInit() {
	staticData := &OutboundAPIParserParserStaticData
	staticData.LiteralNames = []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "'='", "'>'", "'<'", "'<='", "'>='", "'+'", "'-'",
		"'*'", "", "", "", "", "", "'('", "')'", "'['", "']'", "", "", "", "",
		"'\"'", "':'", "','", "';'", "'.'",
	}
	staticData.SymbolicNames = []string{
		"", "CI_VALIDITY_PERIOD_START", "CI_VALIDITY_PERIOD_END", "CI_INSTANCE_CODE",
		"CI_BIDID", "CI_ISP", "CI_DIRECTION", "CI_STATUS", "CI_OPTION_EXPIRY",
		"TIME_INTERVAL_FUNCTION_NAME", "TIME_INTERVAL_GAS_FUNCTION_NAME", "TIME_INTERVAL_EXPLICIT_FUNCTION_NAME",
		"LATEST_GLOBAL", "POINT_IN_TIME_FUNCTION_NAME", "POINT_IN_TIME_UTC_FUNCTION_NAME",
		"TIME_INTERVAL_TO_POINT_IN_TIME_FUNCTION", "RANK_OVER_FUNCTION_NAME",
		"LATEST_FUNCTION_NAME", "IN", "COMPARISON_OPERATOR", "SORT_ORDER", "OPEN_FILTER_INTERVAL_MARKER",
		"EQUAL", "GT", "LT", "LE", "GE", "ADD", "SUB", "MUL", "TIME_ZONE_IANA",
		"DATE", "TIME", "POINT_IN_TIME", "TIME_PERIOD", "LB", "RB", "LSB", "RSB",
		"ID", "SIGNED_INTEGER", "FLOAT", "WORD", "QUOTE", "COLON", "COMMA",
		"SEMICOLON", "DECIMAL_POINT", "WS", "ERRORCHAR",
	}
	staticData.RuleNames = []string{
		"expressionsSection", "keyFilterSection", "keyComparison", "keySurfaceColumn",
		"textColumn", "latestGlobalFunction", "timeInterval", "timeIntervalOrFunction",
		"gasIntervalOrFunction", "pointInTimeOrFunction", "pointInTimeArithmetic",
		"timeIntervalArithmetic", "timeIntervalToPointInTime", "rankOverFunction",
		"rankOverFilter", "latestFunction", "latestExpression", "genericValue",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 49, 236, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 1, 0, 4, 0, 38, 8, 0, 11, 0, 12, 0, 39, 1,
		1, 1, 1, 1, 1, 5, 1, 45, 8, 1, 10, 1, 12, 1, 48, 9, 1, 1, 1, 3, 1, 51,
		8, 1, 1, 1, 1, 1, 1, 2, 1, 2, 3, 2, 57, 8, 2, 1, 2, 1, 2, 1, 2, 1, 2, 3,
		2, 63, 8, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2,
		3, 2, 75, 8, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1,
		2, 3, 2, 87, 8, 2, 1, 3, 1, 3, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 6,
		1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7,
		1, 7, 3, 7, 111, 8, 7, 1, 7, 1, 7, 3, 7, 115, 8, 7, 1, 8, 1, 8, 1, 8, 1,
		8, 1, 8, 3, 8, 122, 8, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1,
		9, 1, 9, 1, 9, 1, 9, 1, 9, 3, 9, 136, 8, 9, 1, 9, 1, 9, 1, 9, 3, 9, 141,
		8, 9, 1, 10, 1, 10, 1, 10, 3, 10, 146, 8, 10, 1, 11, 1, 11, 1, 11, 3, 11,
		151, 8, 11, 1, 11, 3, 11, 154, 8, 11, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12,
		1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 5, 13, 167, 8, 13, 10, 13, 12,
		13, 170, 9, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13,
		5, 13, 180, 8, 13, 10, 13, 12, 13, 183, 9, 13, 1, 13, 1, 13, 1, 13, 5,
		13, 188, 8, 13, 10, 13, 12, 13, 191, 9, 13, 1, 13, 1, 13, 1, 14, 1, 14,
		1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 3, 14, 206,
		8, 14, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 5, 15, 213, 8, 15, 10, 15, 12,
		15, 216, 9, 15, 1, 15, 1, 15, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16,
		1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 3, 16, 232, 8, 16, 1, 17, 1,
		17, 1, 17, 0, 0, 18, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26,
		28, 30, 32, 34, 0, 5, 1, 0, 40, 41, 2, 0, 1, 3, 5, 5, 2, 0, 4, 4, 6, 8,
		1, 0, 27, 28, 4, 0, 18, 18, 20, 21, 30, 34, 39, 42, 248, 0, 37, 1, 0, 0,
		0, 2, 41, 1, 0, 0, 0, 4, 86, 1, 0, 0, 0, 6, 88, 1, 0, 0, 0, 8, 90, 1, 0,
		0, 0, 10, 92, 1, 0, 0, 0, 12, 96, 1, 0, 0, 0, 14, 114, 1, 0, 0, 0, 16,
		116, 1, 0, 0, 0, 18, 140, 1, 0, 0, 0, 20, 142, 1, 0, 0, 0, 22, 153, 1,
		0, 0, 0, 24, 155, 1, 0, 0, 0, 26, 160, 1, 0, 0, 0, 28, 205, 1, 0, 0, 0,
		30, 207, 1, 0, 0, 0, 32, 231, 1, 0, 0, 0, 34, 233, 1, 0, 0, 0, 36, 38,
		3, 2, 1, 0, 37, 36, 1, 0, 0, 0, 38, 39, 1, 0, 0, 0, 39, 37, 1, 0, 0, 0,
		39, 40, 1, 0, 0, 0, 40, 1, 1, 0, 0, 0, 41, 46, 3, 4, 2, 0, 42, 43, 5, 46,
		0, 0, 43, 45, 3, 4, 2, 0, 44, 42, 1, 0, 0, 0, 45, 48, 1, 0, 0, 0, 46, 44,
		1, 0, 0, 0, 46, 47, 1, 0, 0, 0, 47, 50, 1, 0, 0, 0, 48, 46, 1, 0, 0, 0,
		49, 51, 5, 46, 0, 0, 50, 49, 1, 0, 0, 0, 50, 51, 1, 0, 0, 0, 51, 52, 1,
		0, 0, 0, 52, 53, 5, 0, 0, 1, 53, 3, 1, 0, 0, 0, 54, 57, 5, 39, 0, 0, 55,
		57, 3, 6, 3, 0, 56, 54, 1, 0, 0, 0, 56, 55, 1, 0, 0, 0, 57, 58, 1, 0, 0,
		0, 58, 59, 5, 19, 0, 0, 59, 87, 3, 20, 10, 0, 60, 63, 5, 39, 0, 0, 61,
		63, 3, 6, 3, 0, 62, 60, 1, 0, 0, 0, 62, 61, 1, 0, 0, 0, 63, 64, 1, 0, 0,
		0, 64, 65, 5, 18, 0, 0, 65, 87, 3, 22, 11, 0, 66, 67, 5, 39, 0, 0, 67,
		68, 5, 19, 0, 0, 68, 87, 7, 0, 0, 0, 69, 70, 5, 39, 0, 0, 70, 71, 5, 19,
		0, 0, 71, 87, 3, 10, 5, 0, 72, 75, 5, 39, 0, 0, 73, 75, 3, 6, 3, 0, 74,
		72, 1, 0, 0, 0, 74, 73, 1, 0, 0, 0, 75, 76, 1, 0, 0, 0, 76, 77, 5, 19,
		0, 0, 77, 87, 3, 24, 12, 0, 78, 79, 5, 39, 0, 0, 79, 80, 5, 19, 0, 0, 80,
		87, 3, 30, 15, 0, 81, 82, 3, 8, 4, 0, 82, 83, 5, 19, 0, 0, 83, 84, 3, 34,
		17, 0, 84, 87, 1, 0, 0, 0, 85, 87, 3, 26, 13, 0, 86, 56, 1, 0, 0, 0, 86,
		62, 1, 0, 0, 0, 86, 66, 1, 0, 0, 0, 86, 69, 1, 0, 0, 0, 86, 74, 1, 0, 0,
		0, 86, 78, 1, 0, 0, 0, 86, 81, 1, 0, 0, 0, 86, 85, 1, 0, 0, 0, 87, 5, 1,
		0, 0, 0, 88, 89, 7, 1, 0, 0, 89, 7, 1, 0, 0, 0, 90, 91, 7, 2, 0, 0, 91,
		9, 1, 0, 0, 0, 92, 93, 5, 12, 0, 0, 93, 94, 5, 35, 0, 0, 94, 95, 5, 36,
		0, 0, 95, 11, 1, 0, 0, 0, 96, 97, 5, 33, 0, 0, 97, 98, 5, 45, 0, 0, 98,
		99, 5, 33, 0, 0, 99, 13, 1, 0, 0, 0, 100, 101, 5, 11, 0, 0, 101, 102, 5,
		35, 0, 0, 102, 103, 3, 12, 6, 0, 103, 104, 5, 36, 0, 0, 104, 115, 1, 0,
		0, 0, 105, 106, 5, 9, 0, 0, 106, 107, 5, 35, 0, 0, 107, 110, 3, 20, 10,
		0, 108, 109, 5, 45, 0, 0, 109, 111, 5, 30, 0, 0, 110, 108, 1, 0, 0, 0,
		110, 111, 1, 0, 0, 0, 111, 112, 1, 0, 0, 0, 112, 113, 5, 36, 0, 0, 113,
		115, 1, 0, 0, 0, 114, 100, 1, 0, 0, 0, 114, 105, 1, 0, 0, 0, 115, 15, 1,
		0, 0, 0, 116, 117, 5, 10, 0, 0, 117, 118, 5, 35, 0, 0, 118, 121, 3, 20,
		10, 0, 119, 120, 5, 45, 0, 0, 120, 122, 5, 30, 0, 0, 121, 119, 1, 0, 0,
		0, 121, 122, 1, 0, 0, 0, 122, 123, 1, 0, 0, 0, 123, 124, 5, 36, 0, 0, 124,
		17, 1, 0, 0, 0, 125, 141, 5, 33, 0, 0, 126, 127, 5, 13, 0, 0, 127, 128,
		5, 35, 0, 0, 128, 141, 5, 36, 0, 0, 129, 130, 5, 14, 0, 0, 130, 135, 5,
		35, 0, 0, 131, 136, 5, 33, 0, 0, 132, 133, 5, 13, 0, 0, 133, 134, 5, 35,
		0, 0, 134, 136, 5, 36, 0, 0, 135, 131, 1, 0, 0, 0, 135, 132, 1, 0, 0, 0,
		136, 137, 1, 0, 0, 0, 137, 138, 5, 45, 0, 0, 138, 139, 5, 30, 0, 0, 139,
		141, 5, 36, 0, 0, 140, 125, 1, 0, 0, 0, 140, 126, 1, 0, 0, 0, 140, 129,
		1, 0, 0, 0, 141, 19, 1, 0, 0, 0, 142, 145, 3, 18, 9, 0, 143, 144, 7, 3,
		0, 0, 144, 146, 5, 34, 0, 0, 145, 143, 1, 0, 0, 0, 145, 146, 1, 0, 0, 0,
		146, 21, 1, 0, 0, 0, 147, 150, 3, 14, 7, 0, 148, 149, 7, 3, 0, 0, 149,
		151, 5, 34, 0, 0, 150, 148, 1, 0, 0, 0, 150, 151, 1, 0, 0, 0, 151, 154,
		1, 0, 0, 0, 152, 154, 3, 16, 8, 0, 153, 147, 1, 0, 0, 0, 153, 152, 1, 0,
		0, 0, 154, 23, 1, 0, 0, 0, 155, 156, 5, 15, 0, 0, 156, 157, 5, 35, 0, 0,
		157, 158, 3, 22, 11, 0, 158, 159, 5, 36, 0, 0, 159, 25, 1, 0, 0, 0, 160,
		161, 5, 16, 0, 0, 161, 162, 5, 35, 0, 0, 162, 163, 5, 37, 0, 0, 163, 168,
		5, 39, 0, 0, 164, 165, 5, 45, 0, 0, 165, 167, 5, 39, 0, 0, 166, 164, 1,
		0, 0, 0, 167, 170, 1, 0, 0, 0, 168, 166, 1, 0, 0, 0, 168, 169, 1, 0, 0,
		0, 169, 171, 1, 0, 0, 0, 170, 168, 1, 0, 0, 0, 171, 172, 5, 38, 0, 0, 172,
		173, 5, 45, 0, 0, 173, 174, 5, 37, 0, 0, 174, 175, 5, 39, 0, 0, 175, 181,
		5, 20, 0, 0, 176, 177, 5, 45, 0, 0, 177, 178, 5, 39, 0, 0, 178, 180, 5,
		20, 0, 0, 179, 176, 1, 0, 0, 0, 180, 183, 1, 0, 0, 0, 181, 179, 1, 0, 0,
		0, 181, 182, 1, 0, 0, 0, 182, 184, 1, 0, 0, 0, 183, 181, 1, 0, 0, 0, 184,
		189, 5, 38, 0, 0, 185, 186, 5, 45, 0, 0, 186, 188, 3, 28, 14, 0, 187, 185,
		1, 0, 0, 0, 188, 191, 1, 0, 0, 0, 189, 187, 1, 0, 0, 0, 189, 190, 1, 0,
		0, 0, 190, 192, 1, 0, 0, 0, 191, 189, 1, 0, 0, 0, 192, 193, 5, 36, 0, 0,
		193, 27, 1, 0, 0, 0, 194, 206, 5, 40, 0, 0, 195, 196, 5, 37, 0, 0, 196,
		197, 5, 40, 0, 0, 197, 198, 5, 45, 0, 0, 198, 199, 5, 40, 0, 0, 199, 206,
		5, 38, 0, 0, 200, 201, 5, 37, 0, 0, 201, 202, 5, 40, 0, 0, 202, 203, 5,
		45, 0, 0, 203, 204, 5, 21, 0, 0, 204, 206, 5, 38, 0, 0, 205, 194, 1, 0,
		0, 0, 205, 195, 1, 0, 0, 0, 205, 200, 1, 0, 0, 0, 206, 29, 1, 0, 0, 0,
		207, 208, 5, 17, 0, 0, 208, 209, 5, 35, 0, 0, 209, 214, 3, 32, 16, 0, 210,
		211, 5, 45, 0, 0, 211, 213, 3, 32, 16, 0, 212, 210, 1, 0, 0, 0, 213, 216,
		1, 0, 0, 0, 214, 212, 1, 0, 0, 0, 214, 215, 1, 0, 0, 0, 215, 217, 1, 0,
		0, 0, 216, 214, 1, 0, 0, 0, 217, 218, 5, 36, 0, 0, 218, 31, 1, 0, 0, 0,
		219, 220, 5, 39, 0, 0, 220, 221, 5, 19, 0, 0, 221, 232, 3, 20, 10, 0, 222,
		223, 5, 39, 0, 0, 223, 224, 5, 18, 0, 0, 224, 232, 3, 22, 11, 0, 225, 226,
		5, 39, 0, 0, 226, 227, 5, 19, 0, 0, 227, 232, 3, 24, 12, 0, 228, 229, 5,
		39, 0, 0, 229, 230, 5, 19, 0, 0, 230, 232, 7, 0, 0, 0, 231, 219, 1, 0,
		0, 0, 231, 222, 1, 0, 0, 0, 231, 225, 1, 0, 0, 0, 231, 228, 1, 0, 0, 0,
		232, 33, 1, 0, 0, 0, 233, 234, 7, 4, 0, 0, 234, 35, 1, 0, 0, 0, 21, 39,
		46, 50, 56, 62, 74, 86, 110, 114, 121, 135, 140, 145, 150, 153, 168, 181,
		189, 205, 214, 231,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// OutboundAPIParserInit initializes any static state used to implement OutboundAPIParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewOutboundAPIParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func OutboundAPIParserInit() {
	staticData := &OutboundAPIParserParserStaticData
	staticData.once.Do(outboundapiparserParserInit)
}

// NewOutboundAPIParser produces a new parser instance for the optional input antlr.TokenStream.
func NewOutboundAPIParser(input antlr.TokenStream) *OutboundAPIParser {
	OutboundAPIParserInit()
	this := new(OutboundAPIParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &OutboundAPIParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "OutboundAPIParser.g4"

	return this
}

// OutboundAPIParser tokens.
const (
	OutboundAPIParserEOF                                     = antlr.TokenEOF
	OutboundAPIParserCI_VALIDITY_PERIOD_START                = 1
	OutboundAPIParserCI_VALIDITY_PERIOD_END                  = 2
	OutboundAPIParserCI_INSTANCE_CODE                        = 3
	OutboundAPIParserCI_BIDID                                = 4
	OutboundAPIParserCI_ISP                                  = 5
	OutboundAPIParserCI_DIRECTION                            = 6
	OutboundAPIParserCI_STATUS                               = 7
	OutboundAPIParserCI_OPTION_EXPIRY                        = 8
	OutboundAPIParserTIME_INTERVAL_FUNCTION_NAME             = 9
	OutboundAPIParserTIME_INTERVAL_GAS_FUNCTION_NAME         = 10
	OutboundAPIParserTIME_INTERVAL_EXPLICIT_FUNCTION_NAME    = 11
	OutboundAPIParserLATEST_GLOBAL                           = 12
	OutboundAPIParserPOINT_IN_TIME_FUNCTION_NAME             = 13
	OutboundAPIParserPOINT_IN_TIME_UTC_FUNCTION_NAME         = 14
	OutboundAPIParserTIME_INTERVAL_TO_POINT_IN_TIME_FUNCTION = 15
	OutboundAPIParserRANK_OVER_FUNCTION_NAME                 = 16
	OutboundAPIParserLATEST_FUNCTION_NAME                    = 17
	OutboundAPIParserIN                                      = 18
	OutboundAPIParserCOMPARISON_OPERATOR                     = 19
	OutboundAPIParserSORT_ORDER                              = 20
	OutboundAPIParserOPEN_FILTER_INTERVAL_MARKER             = 21
	OutboundAPIParserEQUAL                                   = 22
	OutboundAPIParserGT                                      = 23
	OutboundAPIParserLT                                      = 24
	OutboundAPIParserLE                                      = 25
	OutboundAPIParserGE                                      = 26
	OutboundAPIParserADD                                     = 27
	OutboundAPIParserSUB                                     = 28
	OutboundAPIParserMUL                                     = 29
	OutboundAPIParserTIME_ZONE_IANA                          = 30
	OutboundAPIParserDATE                                    = 31
	OutboundAPIParserTIME                                    = 32
	OutboundAPIParserPOINT_IN_TIME                           = 33
	OutboundAPIParserTIME_PERIOD                             = 34
	OutboundAPIParserLB                                      = 35
	OutboundAPIParserRB                                      = 36
	OutboundAPIParserLSB                                     = 37
	OutboundAPIParserRSB                                     = 38
	OutboundAPIParserID                                      = 39
	OutboundAPIParserSIGNED_INTEGER                          = 40
	OutboundAPIParserFLOAT                                   = 41
	OutboundAPIParserWORD                                    = 42
	OutboundAPIParserQUOTE                                   = 43
	OutboundAPIParserCOLON                                   = 44
	OutboundAPIParserCOMMA                                   = 45
	OutboundAPIParserSEMICOLON                               = 46
	OutboundAPIParserDECIMAL_POINT                           = 47
	OutboundAPIParserWS                                      = 48
	OutboundAPIParserERRORCHAR                               = 49
)

// OutboundAPIParser rules.
const (
	OutboundAPIParserRULE_expressionsSection        = 0
	OutboundAPIParserRULE_keyFilterSection          = 1
	OutboundAPIParserRULE_keyComparison             = 2
	OutboundAPIParserRULE_keySurfaceColumn          = 3
	OutboundAPIParserRULE_textColumn                = 4
	OutboundAPIParserRULE_latestGlobalFunction      = 5
	OutboundAPIParserRULE_timeInterval              = 6
	OutboundAPIParserRULE_timeIntervalOrFunction    = 7
	OutboundAPIParserRULE_gasIntervalOrFunction     = 8
	OutboundAPIParserRULE_pointInTimeOrFunction     = 9
	OutboundAPIParserRULE_pointInTimeArithmetic     = 10
	OutboundAPIParserRULE_timeIntervalArithmetic    = 11
	OutboundAPIParserRULE_timeIntervalToPointInTime = 12
	OutboundAPIParserRULE_rankOverFunction          = 13
	OutboundAPIParserRULE_rankOverFilter            = 14
	OutboundAPIParserRULE_latestFunction            = 15
	OutboundAPIParserRULE_latestExpression          = 16
	OutboundAPIParserRULE_genericValue              = 17
)

// IExpressionsSectionContext is an interface to support dynamic dispatch.
type IExpressionsSectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllKeyFilterSection() []IKeyFilterSectionContext
	KeyFilterSection(i int) IKeyFilterSectionContext

	// IsExpressionsSectionContext differentiates from other interfaces.
	IsExpressionsSectionContext()
}

type ExpressionsSectionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionsSectionContext() *ExpressionsSectionContext {
	var p = new(ExpressionsSectionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_expressionsSection
	return p
}

func InitEmptyExpressionsSectionContext(p *ExpressionsSectionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_expressionsSection
}

func (*ExpressionsSectionContext) IsExpressionsSectionContext() {}

func NewExpressionsSectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionsSectionContext {
	var p = new(ExpressionsSectionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_expressionsSection

	return p
}

func (s *ExpressionsSectionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionsSectionContext) AllKeyFilterSection() []IKeyFilterSectionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IKeyFilterSectionContext); ok {
			len++
		}
	}

	tst := make([]IKeyFilterSectionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IKeyFilterSectionContext); ok {
			tst[i] = t.(IKeyFilterSectionContext)
			i++
		}
	}

	return tst
}

func (s *ExpressionsSectionContext) KeyFilterSection(i int) IKeyFilterSectionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IKeyFilterSectionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IKeyFilterSectionContext)
}

func (s *ExpressionsSectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionsSectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionsSectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterExpressionsSection(s)
	}
}

func (s *ExpressionsSectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitExpressionsSection(s)
	}
}

func (s *ExpressionsSectionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitExpressionsSection(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) ExpressionsSection() (localctx IExpressionsSectionContext) {
	localctx = NewExpressionsSectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, OutboundAPIParserRULE_expressionsSection)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(37)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&549755879934) != 0) {
		{
			p.SetState(36)
			p.KeyFilterSection()
		}

		p.SetState(39)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IKeyFilterSectionContext is an interface to support dynamic dispatch.
type IKeyFilterSectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllKeyComparison() []IKeyComparisonContext
	KeyComparison(i int) IKeyComparisonContext
	EOF() antlr.TerminalNode
	AllSEMICOLON() []antlr.TerminalNode
	SEMICOLON(i int) antlr.TerminalNode

	// IsKeyFilterSectionContext differentiates from other interfaces.
	IsKeyFilterSectionContext()
}

type KeyFilterSectionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyKeyFilterSectionContext() *KeyFilterSectionContext {
	var p = new(KeyFilterSectionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_keyFilterSection
	return p
}

func InitEmptyKeyFilterSectionContext(p *KeyFilterSectionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_keyFilterSection
}

func (*KeyFilterSectionContext) IsKeyFilterSectionContext() {}

func NewKeyFilterSectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *KeyFilterSectionContext {
	var p = new(KeyFilterSectionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_keyFilterSection

	return p
}

func (s *KeyFilterSectionContext) GetParser() antlr.Parser { return s.parser }

func (s *KeyFilterSectionContext) AllKeyComparison() []IKeyComparisonContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IKeyComparisonContext); ok {
			len++
		}
	}

	tst := make([]IKeyComparisonContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IKeyComparisonContext); ok {
			tst[i] = t.(IKeyComparisonContext)
			i++
		}
	}

	return tst
}

func (s *KeyFilterSectionContext) KeyComparison(i int) IKeyComparisonContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IKeyComparisonContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IKeyComparisonContext)
}

func (s *KeyFilterSectionContext) EOF() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserEOF, 0)
}

func (s *KeyFilterSectionContext) AllSEMICOLON() []antlr.TerminalNode {
	return s.GetTokens(OutboundAPIParserSEMICOLON)
}

func (s *KeyFilterSectionContext) SEMICOLON(i int) antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserSEMICOLON, i)
}

func (s *KeyFilterSectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *KeyFilterSectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *KeyFilterSectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterKeyFilterSection(s)
	}
}

func (s *KeyFilterSectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitKeyFilterSection(s)
	}
}

func (s *KeyFilterSectionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitKeyFilterSection(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) KeyFilterSection() (localctx IKeyFilterSectionContext) {
	localctx = NewKeyFilterSectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, OutboundAPIParserRULE_keyFilterSection)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(41)
		p.KeyComparison()
	}
	p.SetState(46)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(42)
				p.Match(OutboundAPIParserSEMICOLON)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(43)
				p.KeyComparison()
			}

		}
		p.SetState(48)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	p.SetState(50)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == OutboundAPIParserSEMICOLON {
		{
			p.SetState(49)
			p.Match(OutboundAPIParserSEMICOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(52)
		p.Match(OutboundAPIParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IKeyComparisonContext is an interface to support dynamic dispatch.
type IKeyComparisonContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsKeyComparisonContext differentiates from other interfaces.
	IsKeyComparisonContext()
}

type KeyComparisonContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyKeyComparisonContext() *KeyComparisonContext {
	var p = new(KeyComparisonContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_keyComparison
	return p
}

func InitEmptyKeyComparisonContext(p *KeyComparisonContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_keyComparison
}

func (*KeyComparisonContext) IsKeyComparisonContext() {}

func NewKeyComparisonContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *KeyComparisonContext {
	var p = new(KeyComparisonContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_keyComparison

	return p
}

func (s *KeyComparisonContext) GetParser() antlr.Parser { return s.parser }

func (s *KeyComparisonContext) CopyAll(ctx *KeyComparisonContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *KeyComparisonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *KeyComparisonContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type IdNumericComparisonContext struct {
	KeyComparisonContext
	number antlr.Token
}

func NewIdNumericComparisonContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IdNumericComparisonContext {
	var p = new(IdNumericComparisonContext)

	InitEmptyKeyComparisonContext(&p.KeyComparisonContext)
	p.parser = parser
	p.CopyAll(ctx.(*KeyComparisonContext))

	return p
}

func (s *IdNumericComparisonContext) GetNumber() antlr.Token { return s.number }

func (s *IdNumericComparisonContext) SetNumber(v antlr.Token) { s.number = v }

func (s *IdNumericComparisonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdNumericComparisonContext) ID() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserID, 0)
}

func (s *IdNumericComparisonContext) COMPARISON_OPERATOR() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCOMPARISON_OPERATOR, 0)
}

func (s *IdNumericComparisonContext) SIGNED_INTEGER() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserSIGNED_INTEGER, 0)
}

func (s *IdNumericComparisonContext) FLOAT() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserFLOAT, 0)
}

func (s *IdNumericComparisonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterIdNumericComparison(s)
	}
}

func (s *IdNumericComparisonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitIdNumericComparison(s)
	}
}

func (s *IdNumericComparisonContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitIdNumericComparison(s)

	default:
		return t.VisitChildren(s)
	}
}

type IdLatestGlobalComparisonContext struct {
	KeyComparisonContext
}

func NewIdLatestGlobalComparisonContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IdLatestGlobalComparisonContext {
	var p = new(IdLatestGlobalComparisonContext)

	InitEmptyKeyComparisonContext(&p.KeyComparisonContext)
	p.parser = parser
	p.CopyAll(ctx.(*KeyComparisonContext))

	return p
}

func (s *IdLatestGlobalComparisonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdLatestGlobalComparisonContext) ID() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserID, 0)
}

func (s *IdLatestGlobalComparisonContext) COMPARISON_OPERATOR() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCOMPARISON_OPERATOR, 0)
}

func (s *IdLatestGlobalComparisonContext) LatestGlobalFunction() ILatestGlobalFunctionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILatestGlobalFunctionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILatestGlobalFunctionContext)
}

func (s *IdLatestGlobalComparisonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterIdLatestGlobalComparison(s)
	}
}

func (s *IdLatestGlobalComparisonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitIdLatestGlobalComparison(s)
	}
}

func (s *IdLatestGlobalComparisonContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitIdLatestGlobalComparison(s)

	default:
		return t.VisitChildren(s)
	}
}

type IdTimeIntervalToPointInTimeComparisonContext struct {
	KeyComparisonContext
}

func NewIdTimeIntervalToPointInTimeComparisonContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IdTimeIntervalToPointInTimeComparisonContext {
	var p = new(IdTimeIntervalToPointInTimeComparisonContext)

	InitEmptyKeyComparisonContext(&p.KeyComparisonContext)
	p.parser = parser
	p.CopyAll(ctx.(*KeyComparisonContext))

	return p
}

func (s *IdTimeIntervalToPointInTimeComparisonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdTimeIntervalToPointInTimeComparisonContext) COMPARISON_OPERATOR() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCOMPARISON_OPERATOR, 0)
}

func (s *IdTimeIntervalToPointInTimeComparisonContext) TimeIntervalToPointInTime() ITimeIntervalToPointInTimeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITimeIntervalToPointInTimeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITimeIntervalToPointInTimeContext)
}

func (s *IdTimeIntervalToPointInTimeComparisonContext) ID() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserID, 0)
}

func (s *IdTimeIntervalToPointInTimeComparisonContext) KeySurfaceColumn() IKeySurfaceColumnContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IKeySurfaceColumnContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IKeySurfaceColumnContext)
}

func (s *IdTimeIntervalToPointInTimeComparisonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterIdTimeIntervalToPointInTimeComparison(s)
	}
}

func (s *IdTimeIntervalToPointInTimeComparisonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitIdTimeIntervalToPointInTimeComparison(s)
	}
}

func (s *IdTimeIntervalToPointInTimeComparisonContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitIdTimeIntervalToPointInTimeComparison(s)

	default:
		return t.VisitChildren(s)
	}
}

type IdTimeIntervalInContext struct {
	KeyComparisonContext
}

func NewIdTimeIntervalInContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IdTimeIntervalInContext {
	var p = new(IdTimeIntervalInContext)

	InitEmptyKeyComparisonContext(&p.KeyComparisonContext)
	p.parser = parser
	p.CopyAll(ctx.(*KeyComparisonContext))

	return p
}

func (s *IdTimeIntervalInContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdTimeIntervalInContext) IN() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserIN, 0)
}

func (s *IdTimeIntervalInContext) TimeIntervalArithmetic() ITimeIntervalArithmeticContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITimeIntervalArithmeticContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITimeIntervalArithmeticContext)
}

func (s *IdTimeIntervalInContext) ID() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserID, 0)
}

func (s *IdTimeIntervalInContext) KeySurfaceColumn() IKeySurfaceColumnContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IKeySurfaceColumnContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IKeySurfaceColumnContext)
}

func (s *IdTimeIntervalInContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterIdTimeIntervalIn(s)
	}
}

func (s *IdTimeIntervalInContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitIdTimeIntervalIn(s)
	}
}

func (s *IdTimeIntervalInContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitIdTimeIntervalIn(s)

	default:
		return t.VisitChildren(s)
	}
}

type IdLatestComparisonContext struct {
	KeyComparisonContext
}

func NewIdLatestComparisonContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IdLatestComparisonContext {
	var p = new(IdLatestComparisonContext)

	InitEmptyKeyComparisonContext(&p.KeyComparisonContext)
	p.parser = parser
	p.CopyAll(ctx.(*KeyComparisonContext))

	return p
}

func (s *IdLatestComparisonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdLatestComparisonContext) ID() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserID, 0)
}

func (s *IdLatestComparisonContext) COMPARISON_OPERATOR() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCOMPARISON_OPERATOR, 0)
}

func (s *IdLatestComparisonContext) LatestFunction() ILatestFunctionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILatestFunctionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILatestFunctionContext)
}

func (s *IdLatestComparisonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterIdLatestComparison(s)
	}
}

func (s *IdLatestComparisonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitIdLatestComparison(s)
	}
}

func (s *IdLatestComparisonContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitIdLatestComparison(s)

	default:
		return t.VisitChildren(s)
	}
}

type RankOverContext struct {
	KeyComparisonContext
}

func NewRankOverContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *RankOverContext {
	var p = new(RankOverContext)

	InitEmptyKeyComparisonContext(&p.KeyComparisonContext)
	p.parser = parser
	p.CopyAll(ctx.(*KeyComparisonContext))

	return p
}

func (s *RankOverContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RankOverContext) RankOverFunction() IRankOverFunctionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRankOverFunctionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRankOverFunctionContext)
}

func (s *RankOverContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterRankOver(s)
	}
}

func (s *RankOverContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitRankOver(s)
	}
}

func (s *RankOverContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitRankOver(s)

	default:
		return t.VisitChildren(s)
	}
}

type IdPointInTimeArithmeticComparisonContext struct {
	KeyComparisonContext
}

func NewIdPointInTimeArithmeticComparisonContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IdPointInTimeArithmeticComparisonContext {
	var p = new(IdPointInTimeArithmeticComparisonContext)

	InitEmptyKeyComparisonContext(&p.KeyComparisonContext)
	p.parser = parser
	p.CopyAll(ctx.(*KeyComparisonContext))

	return p
}

func (s *IdPointInTimeArithmeticComparisonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdPointInTimeArithmeticComparisonContext) COMPARISON_OPERATOR() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCOMPARISON_OPERATOR, 0)
}

func (s *IdPointInTimeArithmeticComparisonContext) PointInTimeArithmetic() IPointInTimeArithmeticContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPointInTimeArithmeticContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPointInTimeArithmeticContext)
}

func (s *IdPointInTimeArithmeticComparisonContext) ID() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserID, 0)
}

func (s *IdPointInTimeArithmeticComparisonContext) KeySurfaceColumn() IKeySurfaceColumnContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IKeySurfaceColumnContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IKeySurfaceColumnContext)
}

func (s *IdPointInTimeArithmeticComparisonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterIdPointInTimeArithmeticComparison(s)
	}
}

func (s *IdPointInTimeArithmeticComparisonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitIdPointInTimeArithmeticComparison(s)
	}
}

func (s *IdPointInTimeArithmeticComparisonContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitIdPointInTimeArithmeticComparison(s)

	default:
		return t.VisitChildren(s)
	}
}

type TextComparisonContext struct {
	KeyComparisonContext
}

func NewTextComparisonContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TextComparisonContext {
	var p = new(TextComparisonContext)

	InitEmptyKeyComparisonContext(&p.KeyComparisonContext)
	p.parser = parser
	p.CopyAll(ctx.(*KeyComparisonContext))

	return p
}

func (s *TextComparisonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TextComparisonContext) TextColumn() ITextColumnContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITextColumnContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITextColumnContext)
}

func (s *TextComparisonContext) COMPARISON_OPERATOR() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCOMPARISON_OPERATOR, 0)
}

func (s *TextComparisonContext) GenericValue() IGenericValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGenericValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGenericValueContext)
}

func (s *TextComparisonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterTextComparison(s)
	}
}

func (s *TextComparisonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitTextComparison(s)
	}
}

func (s *TextComparisonContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitTextComparison(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) KeyComparison() (localctx IKeyComparisonContext) {
	localctx = NewKeyComparisonContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, OutboundAPIParserRULE_keyComparison)
	var _la int

	p.SetState(86)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 6, p.GetParserRuleContext()) {
	case 1:
		localctx = NewIdPointInTimeArithmeticComparisonContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		p.SetState(56)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case OutboundAPIParserID:
			{
				p.SetState(54)
				p.Match(OutboundAPIParserID)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case OutboundAPIParserCI_VALIDITY_PERIOD_START, OutboundAPIParserCI_VALIDITY_PERIOD_END, OutboundAPIParserCI_INSTANCE_CODE, OutboundAPIParserCI_ISP:
			{
				p.SetState(55)
				p.KeySurfaceColumn()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}
		{
			p.SetState(58)
			p.Match(OutboundAPIParserCOMPARISON_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(59)
			p.PointInTimeArithmetic()
		}

	case 2:
		localctx = NewIdTimeIntervalInContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		p.SetState(62)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case OutboundAPIParserID:
			{
				p.SetState(60)
				p.Match(OutboundAPIParserID)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case OutboundAPIParserCI_VALIDITY_PERIOD_START, OutboundAPIParserCI_VALIDITY_PERIOD_END, OutboundAPIParserCI_INSTANCE_CODE, OutboundAPIParserCI_ISP:
			{
				p.SetState(61)
				p.KeySurfaceColumn()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}
		{
			p.SetState(64)
			p.Match(OutboundAPIParserIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(65)
			p.TimeIntervalArithmetic()
		}

	case 3:
		localctx = NewIdNumericComparisonContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(66)
			p.Match(OutboundAPIParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(67)
			p.Match(OutboundAPIParserCOMPARISON_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(68)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*IdNumericComparisonContext).number = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == OutboundAPIParserSIGNED_INTEGER || _la == OutboundAPIParserFLOAT) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*IdNumericComparisonContext).number = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	case 4:
		localctx = NewIdLatestGlobalComparisonContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(69)
			p.Match(OutboundAPIParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(70)
			p.Match(OutboundAPIParserCOMPARISON_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(71)
			p.LatestGlobalFunction()
		}

	case 5:
		localctx = NewIdTimeIntervalToPointInTimeComparisonContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		p.SetState(74)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case OutboundAPIParserID:
			{
				p.SetState(72)
				p.Match(OutboundAPIParserID)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case OutboundAPIParserCI_VALIDITY_PERIOD_START, OutboundAPIParserCI_VALIDITY_PERIOD_END, OutboundAPIParserCI_INSTANCE_CODE, OutboundAPIParserCI_ISP:
			{
				p.SetState(73)
				p.KeySurfaceColumn()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}
		{
			p.SetState(76)
			p.Match(OutboundAPIParserCOMPARISON_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(77)
			p.TimeIntervalToPointInTime()
		}

	case 6:
		localctx = NewIdLatestComparisonContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(78)
			p.Match(OutboundAPIParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(79)
			p.Match(OutboundAPIParserCOMPARISON_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(80)
			p.LatestFunction()
		}

	case 7:
		localctx = NewTextComparisonContext(p, localctx)
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(81)
			p.TextColumn()
		}
		{
			p.SetState(82)
			p.Match(OutboundAPIParserCOMPARISON_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(83)
			p.GenericValue()
		}

	case 8:
		localctx = NewRankOverContext(p, localctx)
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(85)
			p.RankOverFunction()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IKeySurfaceColumnContext is an interface to support dynamic dispatch.
type IKeySurfaceColumnContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CI_VALIDITY_PERIOD_START() antlr.TerminalNode
	CI_VALIDITY_PERIOD_END() antlr.TerminalNode
	CI_INSTANCE_CODE() antlr.TerminalNode
	CI_ISP() antlr.TerminalNode

	// IsKeySurfaceColumnContext differentiates from other interfaces.
	IsKeySurfaceColumnContext()
}

type KeySurfaceColumnContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyKeySurfaceColumnContext() *KeySurfaceColumnContext {
	var p = new(KeySurfaceColumnContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_keySurfaceColumn
	return p
}

func InitEmptyKeySurfaceColumnContext(p *KeySurfaceColumnContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_keySurfaceColumn
}

func (*KeySurfaceColumnContext) IsKeySurfaceColumnContext() {}

func NewKeySurfaceColumnContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *KeySurfaceColumnContext {
	var p = new(KeySurfaceColumnContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_keySurfaceColumn

	return p
}

func (s *KeySurfaceColumnContext) GetParser() antlr.Parser { return s.parser }

func (s *KeySurfaceColumnContext) CI_VALIDITY_PERIOD_START() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCI_VALIDITY_PERIOD_START, 0)
}

func (s *KeySurfaceColumnContext) CI_VALIDITY_PERIOD_END() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCI_VALIDITY_PERIOD_END, 0)
}

func (s *KeySurfaceColumnContext) CI_INSTANCE_CODE() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCI_INSTANCE_CODE, 0)
}

func (s *KeySurfaceColumnContext) CI_ISP() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCI_ISP, 0)
}

func (s *KeySurfaceColumnContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *KeySurfaceColumnContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *KeySurfaceColumnContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterKeySurfaceColumn(s)
	}
}

func (s *KeySurfaceColumnContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitKeySurfaceColumn(s)
	}
}

func (s *KeySurfaceColumnContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitKeySurfaceColumn(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) KeySurfaceColumn() (localctx IKeySurfaceColumnContext) {
	localctx = NewKeySurfaceColumnContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, OutboundAPIParserRULE_keySurfaceColumn)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(88)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&46) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITextColumnContext is an interface to support dynamic dispatch.
type ITextColumnContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CI_BIDID() antlr.TerminalNode
	CI_DIRECTION() antlr.TerminalNode
	CI_STATUS() antlr.TerminalNode
	CI_OPTION_EXPIRY() antlr.TerminalNode

	// IsTextColumnContext differentiates from other interfaces.
	IsTextColumnContext()
}

type TextColumnContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTextColumnContext() *TextColumnContext {
	var p = new(TextColumnContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_textColumn
	return p
}

func InitEmptyTextColumnContext(p *TextColumnContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_textColumn
}

func (*TextColumnContext) IsTextColumnContext() {}

func NewTextColumnContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TextColumnContext {
	var p = new(TextColumnContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_textColumn

	return p
}

func (s *TextColumnContext) GetParser() antlr.Parser { return s.parser }

func (s *TextColumnContext) CI_BIDID() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCI_BIDID, 0)
}

func (s *TextColumnContext) CI_DIRECTION() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCI_DIRECTION, 0)
}

func (s *TextColumnContext) CI_STATUS() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCI_STATUS, 0)
}

func (s *TextColumnContext) CI_OPTION_EXPIRY() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCI_OPTION_EXPIRY, 0)
}

func (s *TextColumnContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TextColumnContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TextColumnContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterTextColumn(s)
	}
}

func (s *TextColumnContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitTextColumn(s)
	}
}

func (s *TextColumnContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitTextColumn(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) TextColumn() (localctx ITextColumnContext) {
	localctx = NewTextColumnContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, OutboundAPIParserRULE_textColumn)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(90)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&464) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILatestGlobalFunctionContext is an interface to support dynamic dispatch.
type ILatestGlobalFunctionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LATEST_GLOBAL() antlr.TerminalNode
	LB() antlr.TerminalNode
	RB() antlr.TerminalNode

	// IsLatestGlobalFunctionContext differentiates from other interfaces.
	IsLatestGlobalFunctionContext()
}

type LatestGlobalFunctionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLatestGlobalFunctionContext() *LatestGlobalFunctionContext {
	var p = new(LatestGlobalFunctionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_latestGlobalFunction
	return p
}

func InitEmptyLatestGlobalFunctionContext(p *LatestGlobalFunctionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_latestGlobalFunction
}

func (*LatestGlobalFunctionContext) IsLatestGlobalFunctionContext() {}

func NewLatestGlobalFunctionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LatestGlobalFunctionContext {
	var p = new(LatestGlobalFunctionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_latestGlobalFunction

	return p
}

func (s *LatestGlobalFunctionContext) GetParser() antlr.Parser { return s.parser }

func (s *LatestGlobalFunctionContext) LATEST_GLOBAL() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserLATEST_GLOBAL, 0)
}

func (s *LatestGlobalFunctionContext) LB() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserLB, 0)
}

func (s *LatestGlobalFunctionContext) RB() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserRB, 0)
}

func (s *LatestGlobalFunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LatestGlobalFunctionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LatestGlobalFunctionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterLatestGlobalFunction(s)
	}
}

func (s *LatestGlobalFunctionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitLatestGlobalFunction(s)
	}
}

func (s *LatestGlobalFunctionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitLatestGlobalFunction(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) LatestGlobalFunction() (localctx ILatestGlobalFunctionContext) {
	localctx = NewLatestGlobalFunctionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, OutboundAPIParserRULE_latestGlobalFunction)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(92)
		p.Match(OutboundAPIParserLATEST_GLOBAL)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(93)
		p.Match(OutboundAPIParserLB)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(94)
		p.Match(OutboundAPIParserRB)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITimeIntervalContext is an interface to support dynamic dispatch.
type ITimeIntervalContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllPOINT_IN_TIME() []antlr.TerminalNode
	POINT_IN_TIME(i int) antlr.TerminalNode
	COMMA() antlr.TerminalNode

	// IsTimeIntervalContext differentiates from other interfaces.
	IsTimeIntervalContext()
}

type TimeIntervalContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTimeIntervalContext() *TimeIntervalContext {
	var p = new(TimeIntervalContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_timeInterval
	return p
}

func InitEmptyTimeIntervalContext(p *TimeIntervalContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_timeInterval
}

func (*TimeIntervalContext) IsTimeIntervalContext() {}

func NewTimeIntervalContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimeIntervalContext {
	var p = new(TimeIntervalContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_timeInterval

	return p
}

func (s *TimeIntervalContext) GetParser() antlr.Parser { return s.parser }

func (s *TimeIntervalContext) AllPOINT_IN_TIME() []antlr.TerminalNode {
	return s.GetTokens(OutboundAPIParserPOINT_IN_TIME)
}

func (s *TimeIntervalContext) POINT_IN_TIME(i int) antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserPOINT_IN_TIME, i)
}

func (s *TimeIntervalContext) COMMA() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCOMMA, 0)
}

func (s *TimeIntervalContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeIntervalContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TimeIntervalContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterTimeInterval(s)
	}
}

func (s *TimeIntervalContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitTimeInterval(s)
	}
}

func (s *TimeIntervalContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitTimeInterval(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) TimeInterval() (localctx ITimeIntervalContext) {
	localctx = NewTimeIntervalContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, OutboundAPIParserRULE_timeInterval)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(96)
		p.Match(OutboundAPIParserPOINT_IN_TIME)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(97)
		p.Match(OutboundAPIParserCOMMA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(98)
		p.Match(OutboundAPIParserPOINT_IN_TIME)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITimeIntervalOrFunctionContext is an interface to support dynamic dispatch.
type ITimeIntervalOrFunctionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetExpressionTimeZone returns the expressionTimeZone token.
	GetExpressionTimeZone() antlr.Token

	// SetExpressionTimeZone sets the expressionTimeZone token.
	SetExpressionTimeZone(antlr.Token)

	// Getter signatures
	TIME_INTERVAL_EXPLICIT_FUNCTION_NAME() antlr.TerminalNode
	LB() antlr.TerminalNode
	TimeInterval() ITimeIntervalContext
	RB() antlr.TerminalNode
	TIME_INTERVAL_FUNCTION_NAME() antlr.TerminalNode
	PointInTimeArithmetic() IPointInTimeArithmeticContext
	COMMA() antlr.TerminalNode
	TIME_ZONE_IANA() antlr.TerminalNode

	// IsTimeIntervalOrFunctionContext differentiates from other interfaces.
	IsTimeIntervalOrFunctionContext()
}

type TimeIntervalOrFunctionContext struct {
	antlr.BaseParserRuleContext
	parser             antlr.Parser
	expressionTimeZone antlr.Token
}

func NewEmptyTimeIntervalOrFunctionContext() *TimeIntervalOrFunctionContext {
	var p = new(TimeIntervalOrFunctionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_timeIntervalOrFunction
	return p
}

func InitEmptyTimeIntervalOrFunctionContext(p *TimeIntervalOrFunctionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_timeIntervalOrFunction
}

func (*TimeIntervalOrFunctionContext) IsTimeIntervalOrFunctionContext() {}

func NewTimeIntervalOrFunctionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimeIntervalOrFunctionContext {
	var p = new(TimeIntervalOrFunctionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_timeIntervalOrFunction

	return p
}

func (s *TimeIntervalOrFunctionContext) GetParser() antlr.Parser { return s.parser }

func (s *TimeIntervalOrFunctionContext) GetExpressionTimeZone() antlr.Token {
	return s.expressionTimeZone
}

func (s *TimeIntervalOrFunctionContext) SetExpressionTimeZone(v antlr.Token) {
	s.expressionTimeZone = v
}

func (s *TimeIntervalOrFunctionContext) TIME_INTERVAL_EXPLICIT_FUNCTION_NAME() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserTIME_INTERVAL_EXPLICIT_FUNCTION_NAME, 0)
}

func (s *TimeIntervalOrFunctionContext) LB() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserLB, 0)
}

func (s *TimeIntervalOrFunctionContext) TimeInterval() ITimeIntervalContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITimeIntervalContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITimeIntervalContext)
}

func (s *TimeIntervalOrFunctionContext) RB() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserRB, 0)
}

func (s *TimeIntervalOrFunctionContext) TIME_INTERVAL_FUNCTION_NAME() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserTIME_INTERVAL_FUNCTION_NAME, 0)
}

func (s *TimeIntervalOrFunctionContext) PointInTimeArithmetic() IPointInTimeArithmeticContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPointInTimeArithmeticContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPointInTimeArithmeticContext)
}

func (s *TimeIntervalOrFunctionContext) COMMA() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCOMMA, 0)
}

func (s *TimeIntervalOrFunctionContext) TIME_ZONE_IANA() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserTIME_ZONE_IANA, 0)
}

func (s *TimeIntervalOrFunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeIntervalOrFunctionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TimeIntervalOrFunctionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterTimeIntervalOrFunction(s)
	}
}

func (s *TimeIntervalOrFunctionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitTimeIntervalOrFunction(s)
	}
}

func (s *TimeIntervalOrFunctionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitTimeIntervalOrFunction(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) TimeIntervalOrFunction() (localctx ITimeIntervalOrFunctionContext) {
	localctx = NewTimeIntervalOrFunctionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, OutboundAPIParserRULE_timeIntervalOrFunction)
	var _la int

	p.SetState(114)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case OutboundAPIParserTIME_INTERVAL_EXPLICIT_FUNCTION_NAME:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(100)
			p.Match(OutboundAPIParserTIME_INTERVAL_EXPLICIT_FUNCTION_NAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(101)
			p.Match(OutboundAPIParserLB)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(102)
			p.TimeInterval()
		}
		{
			p.SetState(103)
			p.Match(OutboundAPIParserRB)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case OutboundAPIParserTIME_INTERVAL_FUNCTION_NAME:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(105)
			p.Match(OutboundAPIParserTIME_INTERVAL_FUNCTION_NAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(106)
			p.Match(OutboundAPIParserLB)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(107)
			p.PointInTimeArithmetic()
		}
		p.SetState(110)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == OutboundAPIParserCOMMA {
			{
				p.SetState(108)
				p.Match(OutboundAPIParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(109)

				var _m = p.Match(OutboundAPIParserTIME_ZONE_IANA)

				localctx.(*TimeIntervalOrFunctionContext).expressionTimeZone = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(112)
			p.Match(OutboundAPIParserRB)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IGasIntervalOrFunctionContext is an interface to support dynamic dispatch.
type IGasIntervalOrFunctionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetExpressionTimeZone returns the expressionTimeZone token.
	GetExpressionTimeZone() antlr.Token

	// SetExpressionTimeZone sets the expressionTimeZone token.
	SetExpressionTimeZone(antlr.Token)

	// Getter signatures
	TIME_INTERVAL_GAS_FUNCTION_NAME() antlr.TerminalNode
	LB() antlr.TerminalNode
	PointInTimeArithmetic() IPointInTimeArithmeticContext
	RB() antlr.TerminalNode
	COMMA() antlr.TerminalNode
	TIME_ZONE_IANA() antlr.TerminalNode

	// IsGasIntervalOrFunctionContext differentiates from other interfaces.
	IsGasIntervalOrFunctionContext()
}

type GasIntervalOrFunctionContext struct {
	antlr.BaseParserRuleContext
	parser             antlr.Parser
	expressionTimeZone antlr.Token
}

func NewEmptyGasIntervalOrFunctionContext() *GasIntervalOrFunctionContext {
	var p = new(GasIntervalOrFunctionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_gasIntervalOrFunction
	return p
}

func InitEmptyGasIntervalOrFunctionContext(p *GasIntervalOrFunctionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_gasIntervalOrFunction
}

func (*GasIntervalOrFunctionContext) IsGasIntervalOrFunctionContext() {}

func NewGasIntervalOrFunctionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GasIntervalOrFunctionContext {
	var p = new(GasIntervalOrFunctionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_gasIntervalOrFunction

	return p
}

func (s *GasIntervalOrFunctionContext) GetParser() antlr.Parser { return s.parser }

func (s *GasIntervalOrFunctionContext) GetExpressionTimeZone() antlr.Token {
	return s.expressionTimeZone
}

func (s *GasIntervalOrFunctionContext) SetExpressionTimeZone(v antlr.Token) { s.expressionTimeZone = v }

func (s *GasIntervalOrFunctionContext) TIME_INTERVAL_GAS_FUNCTION_NAME() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserTIME_INTERVAL_GAS_FUNCTION_NAME, 0)
}

func (s *GasIntervalOrFunctionContext) LB() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserLB, 0)
}

func (s *GasIntervalOrFunctionContext) PointInTimeArithmetic() IPointInTimeArithmeticContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPointInTimeArithmeticContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPointInTimeArithmeticContext)
}

func (s *GasIntervalOrFunctionContext) RB() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserRB, 0)
}

func (s *GasIntervalOrFunctionContext) COMMA() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCOMMA, 0)
}

func (s *GasIntervalOrFunctionContext) TIME_ZONE_IANA() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserTIME_ZONE_IANA, 0)
}

func (s *GasIntervalOrFunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GasIntervalOrFunctionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GasIntervalOrFunctionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterGasIntervalOrFunction(s)
	}
}

func (s *GasIntervalOrFunctionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitGasIntervalOrFunction(s)
	}
}

func (s *GasIntervalOrFunctionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitGasIntervalOrFunction(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) GasIntervalOrFunction() (localctx IGasIntervalOrFunctionContext) {
	localctx = NewGasIntervalOrFunctionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, OutboundAPIParserRULE_gasIntervalOrFunction)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(116)
		p.Match(OutboundAPIParserTIME_INTERVAL_GAS_FUNCTION_NAME)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(117)
		p.Match(OutboundAPIParserLB)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(118)
		p.PointInTimeArithmetic()
	}
	p.SetState(121)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == OutboundAPIParserCOMMA {
		{
			p.SetState(119)
			p.Match(OutboundAPIParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(120)

			var _m = p.Match(OutboundAPIParserTIME_ZONE_IANA)

			localctx.(*GasIntervalOrFunctionContext).expressionTimeZone = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(123)
		p.Match(OutboundAPIParserRB)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPointInTimeOrFunctionContext is an interface to support dynamic dispatch.
type IPointInTimeOrFunctionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetExpressionTimeZone returns the expressionTimeZone token.
	GetExpressionTimeZone() antlr.Token

	// SetExpressionTimeZone sets the expressionTimeZone token.
	SetExpressionTimeZone(antlr.Token)

	// Getter signatures
	POINT_IN_TIME() antlr.TerminalNode
	POINT_IN_TIME_FUNCTION_NAME() antlr.TerminalNode
	AllLB() []antlr.TerminalNode
	LB(i int) antlr.TerminalNode
	AllRB() []antlr.TerminalNode
	RB(i int) antlr.TerminalNode
	POINT_IN_TIME_UTC_FUNCTION_NAME() antlr.TerminalNode
	COMMA() antlr.TerminalNode
	TIME_ZONE_IANA() antlr.TerminalNode

	// IsPointInTimeOrFunctionContext differentiates from other interfaces.
	IsPointInTimeOrFunctionContext()
}

type PointInTimeOrFunctionContext struct {
	antlr.BaseParserRuleContext
	parser             antlr.Parser
	expressionTimeZone antlr.Token
}

func NewEmptyPointInTimeOrFunctionContext() *PointInTimeOrFunctionContext {
	var p = new(PointInTimeOrFunctionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_pointInTimeOrFunction
	return p
}

func InitEmptyPointInTimeOrFunctionContext(p *PointInTimeOrFunctionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_pointInTimeOrFunction
}

func (*PointInTimeOrFunctionContext) IsPointInTimeOrFunctionContext() {}

func NewPointInTimeOrFunctionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PointInTimeOrFunctionContext {
	var p = new(PointInTimeOrFunctionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_pointInTimeOrFunction

	return p
}

func (s *PointInTimeOrFunctionContext) GetParser() antlr.Parser { return s.parser }

func (s *PointInTimeOrFunctionContext) GetExpressionTimeZone() antlr.Token {
	return s.expressionTimeZone
}

func (s *PointInTimeOrFunctionContext) SetExpressionTimeZone(v antlr.Token) { s.expressionTimeZone = v }

func (s *PointInTimeOrFunctionContext) POINT_IN_TIME() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserPOINT_IN_TIME, 0)
}

func (s *PointInTimeOrFunctionContext) POINT_IN_TIME_FUNCTION_NAME() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserPOINT_IN_TIME_FUNCTION_NAME, 0)
}

func (s *PointInTimeOrFunctionContext) AllLB() []antlr.TerminalNode {
	return s.GetTokens(OutboundAPIParserLB)
}

func (s *PointInTimeOrFunctionContext) LB(i int) antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserLB, i)
}

func (s *PointInTimeOrFunctionContext) AllRB() []antlr.TerminalNode {
	return s.GetTokens(OutboundAPIParserRB)
}

func (s *PointInTimeOrFunctionContext) RB(i int) antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserRB, i)
}

func (s *PointInTimeOrFunctionContext) POINT_IN_TIME_UTC_FUNCTION_NAME() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserPOINT_IN_TIME_UTC_FUNCTION_NAME, 0)
}

func (s *PointInTimeOrFunctionContext) COMMA() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCOMMA, 0)
}

func (s *PointInTimeOrFunctionContext) TIME_ZONE_IANA() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserTIME_ZONE_IANA, 0)
}

func (s *PointInTimeOrFunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PointInTimeOrFunctionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PointInTimeOrFunctionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterPointInTimeOrFunction(s)
	}
}

func (s *PointInTimeOrFunctionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitPointInTimeOrFunction(s)
	}
}

func (s *PointInTimeOrFunctionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitPointInTimeOrFunction(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) PointInTimeOrFunction() (localctx IPointInTimeOrFunctionContext) {
	localctx = NewPointInTimeOrFunctionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, OutboundAPIParserRULE_pointInTimeOrFunction)
	p.SetState(140)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case OutboundAPIParserPOINT_IN_TIME:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(125)
			p.Match(OutboundAPIParserPOINT_IN_TIME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case OutboundAPIParserPOINT_IN_TIME_FUNCTION_NAME:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(126)
			p.Match(OutboundAPIParserPOINT_IN_TIME_FUNCTION_NAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(127)
			p.Match(OutboundAPIParserLB)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(128)
			p.Match(OutboundAPIParserRB)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case OutboundAPIParserPOINT_IN_TIME_UTC_FUNCTION_NAME:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(129)
			p.Match(OutboundAPIParserPOINT_IN_TIME_UTC_FUNCTION_NAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(130)
			p.Match(OutboundAPIParserLB)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(135)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case OutboundAPIParserPOINT_IN_TIME:
			{
				p.SetState(131)
				p.Match(OutboundAPIParserPOINT_IN_TIME)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case OutboundAPIParserPOINT_IN_TIME_FUNCTION_NAME:
			{
				p.SetState(132)
				p.Match(OutboundAPIParserPOINT_IN_TIME_FUNCTION_NAME)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(133)
				p.Match(OutboundAPIParserLB)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(134)
				p.Match(OutboundAPIParserRB)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}
		{
			p.SetState(137)
			p.Match(OutboundAPIParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(138)

			var _m = p.Match(OutboundAPIParserTIME_ZONE_IANA)

			localctx.(*PointInTimeOrFunctionContext).expressionTimeZone = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(139)
			p.Match(OutboundAPIParserRB)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPointInTimeArithmeticContext is an interface to support dynamic dispatch.
type IPointInTimeArithmeticContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetArithmeticOperator returns the ArithmeticOperator token.
	GetArithmeticOperator() antlr.Token

	// GetTimePeriod returns the TimePeriod token.
	GetTimePeriod() antlr.Token

	// SetArithmeticOperator sets the ArithmeticOperator token.
	SetArithmeticOperator(antlr.Token)

	// SetTimePeriod sets the TimePeriod token.
	SetTimePeriod(antlr.Token)

	// Getter signatures
	PointInTimeOrFunction() IPointInTimeOrFunctionContext
	TIME_PERIOD() antlr.TerminalNode
	ADD() antlr.TerminalNode
	SUB() antlr.TerminalNode

	// IsPointInTimeArithmeticContext differentiates from other interfaces.
	IsPointInTimeArithmeticContext()
}

type PointInTimeArithmeticContext struct {
	antlr.BaseParserRuleContext
	parser             antlr.Parser
	ArithmeticOperator antlr.Token
	TimePeriod         antlr.Token
}

func NewEmptyPointInTimeArithmeticContext() *PointInTimeArithmeticContext {
	var p = new(PointInTimeArithmeticContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_pointInTimeArithmetic
	return p
}

func InitEmptyPointInTimeArithmeticContext(p *PointInTimeArithmeticContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_pointInTimeArithmetic
}

func (*PointInTimeArithmeticContext) IsPointInTimeArithmeticContext() {}

func NewPointInTimeArithmeticContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PointInTimeArithmeticContext {
	var p = new(PointInTimeArithmeticContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_pointInTimeArithmetic

	return p
}

func (s *PointInTimeArithmeticContext) GetParser() antlr.Parser { return s.parser }

func (s *PointInTimeArithmeticContext) GetArithmeticOperator() antlr.Token {
	return s.ArithmeticOperator
}

func (s *PointInTimeArithmeticContext) GetTimePeriod() antlr.Token { return s.TimePeriod }

func (s *PointInTimeArithmeticContext) SetArithmeticOperator(v antlr.Token) { s.ArithmeticOperator = v }

func (s *PointInTimeArithmeticContext) SetTimePeriod(v antlr.Token) { s.TimePeriod = v }

func (s *PointInTimeArithmeticContext) PointInTimeOrFunction() IPointInTimeOrFunctionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPointInTimeOrFunctionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPointInTimeOrFunctionContext)
}

func (s *PointInTimeArithmeticContext) TIME_PERIOD() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserTIME_PERIOD, 0)
}

func (s *PointInTimeArithmeticContext) ADD() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserADD, 0)
}

func (s *PointInTimeArithmeticContext) SUB() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserSUB, 0)
}

func (s *PointInTimeArithmeticContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PointInTimeArithmeticContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PointInTimeArithmeticContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterPointInTimeArithmetic(s)
	}
}

func (s *PointInTimeArithmeticContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitPointInTimeArithmetic(s)
	}
}

func (s *PointInTimeArithmeticContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitPointInTimeArithmetic(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) PointInTimeArithmetic() (localctx IPointInTimeArithmeticContext) {
	localctx = NewPointInTimeArithmeticContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, OutboundAPIParserRULE_pointInTimeArithmetic)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(142)
		p.PointInTimeOrFunction()
	}

	p.SetState(145)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == OutboundAPIParserADD || _la == OutboundAPIParserSUB {
		{
			p.SetState(143)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*PointInTimeArithmeticContext).ArithmeticOperator = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == OutboundAPIParserADD || _la == OutboundAPIParserSUB) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*PointInTimeArithmeticContext).ArithmeticOperator = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(144)

			var _m = p.Match(OutboundAPIParserTIME_PERIOD)

			localctx.(*PointInTimeArithmeticContext).TimePeriod = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITimeIntervalArithmeticContext is an interface to support dynamic dispatch.
type ITimeIntervalArithmeticContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetArithmeticOperator returns the ArithmeticOperator token.
	GetArithmeticOperator() antlr.Token

	// GetTimePeriod returns the TimePeriod token.
	GetTimePeriod() antlr.Token

	// SetArithmeticOperator sets the ArithmeticOperator token.
	SetArithmeticOperator(antlr.Token)

	// SetTimePeriod sets the TimePeriod token.
	SetTimePeriod(antlr.Token)

	// Getter signatures
	TimeIntervalOrFunction() ITimeIntervalOrFunctionContext
	TIME_PERIOD() antlr.TerminalNode
	ADD() antlr.TerminalNode
	SUB() antlr.TerminalNode
	GasIntervalOrFunction() IGasIntervalOrFunctionContext

	// IsTimeIntervalArithmeticContext differentiates from other interfaces.
	IsTimeIntervalArithmeticContext()
}

type TimeIntervalArithmeticContext struct {
	antlr.BaseParserRuleContext
	parser             antlr.Parser
	ArithmeticOperator antlr.Token
	TimePeriod         antlr.Token
}

func NewEmptyTimeIntervalArithmeticContext() *TimeIntervalArithmeticContext {
	var p = new(TimeIntervalArithmeticContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_timeIntervalArithmetic
	return p
}

func InitEmptyTimeIntervalArithmeticContext(p *TimeIntervalArithmeticContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_timeIntervalArithmetic
}

func (*TimeIntervalArithmeticContext) IsTimeIntervalArithmeticContext() {}

func NewTimeIntervalArithmeticContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimeIntervalArithmeticContext {
	var p = new(TimeIntervalArithmeticContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_timeIntervalArithmetic

	return p
}

func (s *TimeIntervalArithmeticContext) GetParser() antlr.Parser { return s.parser }

func (s *TimeIntervalArithmeticContext) GetArithmeticOperator() antlr.Token {
	return s.ArithmeticOperator
}

func (s *TimeIntervalArithmeticContext) GetTimePeriod() antlr.Token { return s.TimePeriod }

func (s *TimeIntervalArithmeticContext) SetArithmeticOperator(v antlr.Token) {
	s.ArithmeticOperator = v
}

func (s *TimeIntervalArithmeticContext) SetTimePeriod(v antlr.Token) { s.TimePeriod = v }

func (s *TimeIntervalArithmeticContext) TimeIntervalOrFunction() ITimeIntervalOrFunctionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITimeIntervalOrFunctionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITimeIntervalOrFunctionContext)
}

func (s *TimeIntervalArithmeticContext) TIME_PERIOD() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserTIME_PERIOD, 0)
}

func (s *TimeIntervalArithmeticContext) ADD() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserADD, 0)
}

func (s *TimeIntervalArithmeticContext) SUB() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserSUB, 0)
}

func (s *TimeIntervalArithmeticContext) GasIntervalOrFunction() IGasIntervalOrFunctionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGasIntervalOrFunctionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGasIntervalOrFunctionContext)
}

func (s *TimeIntervalArithmeticContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeIntervalArithmeticContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TimeIntervalArithmeticContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterTimeIntervalArithmetic(s)
	}
}

func (s *TimeIntervalArithmeticContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitTimeIntervalArithmetic(s)
	}
}

func (s *TimeIntervalArithmeticContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitTimeIntervalArithmetic(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) TimeIntervalArithmetic() (localctx ITimeIntervalArithmeticContext) {
	localctx = NewTimeIntervalArithmeticContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, OutboundAPIParserRULE_timeIntervalArithmetic)
	var _la int

	p.SetState(153)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case OutboundAPIParserTIME_INTERVAL_FUNCTION_NAME, OutboundAPIParserTIME_INTERVAL_EXPLICIT_FUNCTION_NAME:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(147)
			p.TimeIntervalOrFunction()
		}

		p.SetState(150)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == OutboundAPIParserADD || _la == OutboundAPIParserSUB {
			{
				p.SetState(148)

				var _lt = p.GetTokenStream().LT(1)

				localctx.(*TimeIntervalArithmeticContext).ArithmeticOperator = _lt

				_la = p.GetTokenStream().LA(1)

				if !(_la == OutboundAPIParserADD || _la == OutboundAPIParserSUB) {
					var _ri = p.GetErrorHandler().RecoverInline(p)

					localctx.(*TimeIntervalArithmeticContext).ArithmeticOperator = _ri
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(149)

				var _m = p.Match(OutboundAPIParserTIME_PERIOD)

				localctx.(*TimeIntervalArithmeticContext).TimePeriod = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	case OutboundAPIParserTIME_INTERVAL_GAS_FUNCTION_NAME:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(152)
			p.GasIntervalOrFunction()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITimeIntervalToPointInTimeContext is an interface to support dynamic dispatch.
type ITimeIntervalToPointInTimeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TIME_INTERVAL_TO_POINT_IN_TIME_FUNCTION() antlr.TerminalNode
	LB() antlr.TerminalNode
	TimeIntervalArithmetic() ITimeIntervalArithmeticContext
	RB() antlr.TerminalNode

	// IsTimeIntervalToPointInTimeContext differentiates from other interfaces.
	IsTimeIntervalToPointInTimeContext()
}

type TimeIntervalToPointInTimeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTimeIntervalToPointInTimeContext() *TimeIntervalToPointInTimeContext {
	var p = new(TimeIntervalToPointInTimeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_timeIntervalToPointInTime
	return p
}

func InitEmptyTimeIntervalToPointInTimeContext(p *TimeIntervalToPointInTimeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_timeIntervalToPointInTime
}

func (*TimeIntervalToPointInTimeContext) IsTimeIntervalToPointInTimeContext() {}

func NewTimeIntervalToPointInTimeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimeIntervalToPointInTimeContext {
	var p = new(TimeIntervalToPointInTimeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_timeIntervalToPointInTime

	return p
}

func (s *TimeIntervalToPointInTimeContext) GetParser() antlr.Parser { return s.parser }

func (s *TimeIntervalToPointInTimeContext) TIME_INTERVAL_TO_POINT_IN_TIME_FUNCTION() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserTIME_INTERVAL_TO_POINT_IN_TIME_FUNCTION, 0)
}

func (s *TimeIntervalToPointInTimeContext) LB() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserLB, 0)
}

func (s *TimeIntervalToPointInTimeContext) TimeIntervalArithmetic() ITimeIntervalArithmeticContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITimeIntervalArithmeticContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITimeIntervalArithmeticContext)
}

func (s *TimeIntervalToPointInTimeContext) RB() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserRB, 0)
}

func (s *TimeIntervalToPointInTimeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeIntervalToPointInTimeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TimeIntervalToPointInTimeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterTimeIntervalToPointInTime(s)
	}
}

func (s *TimeIntervalToPointInTimeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitTimeIntervalToPointInTime(s)
	}
}

func (s *TimeIntervalToPointInTimeContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitTimeIntervalToPointInTime(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) TimeIntervalToPointInTime() (localctx ITimeIntervalToPointInTimeContext) {
	localctx = NewTimeIntervalToPointInTimeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, OutboundAPIParserRULE_timeIntervalToPointInTime)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(155)
		p.Match(OutboundAPIParserTIME_INTERVAL_TO_POINT_IN_TIME_FUNCTION)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(156)
		p.Match(OutboundAPIParserLB)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(157)
		p.TimeIntervalArithmetic()
	}
	{
		p.SetState(158)
		p.Match(OutboundAPIParserRB)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRankOverFunctionContext is an interface to support dynamic dispatch.
type IRankOverFunctionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	RANK_OVER_FUNCTION_NAME() antlr.TerminalNode
	LB() antlr.TerminalNode
	AllLSB() []antlr.TerminalNode
	LSB(i int) antlr.TerminalNode
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	AllRSB() []antlr.TerminalNode
	RSB(i int) antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	AllSORT_ORDER() []antlr.TerminalNode
	SORT_ORDER(i int) antlr.TerminalNode
	RB() antlr.TerminalNode
	AllRankOverFilter() []IRankOverFilterContext
	RankOverFilter(i int) IRankOverFilterContext

	// IsRankOverFunctionContext differentiates from other interfaces.
	IsRankOverFunctionContext()
}

type RankOverFunctionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRankOverFunctionContext() *RankOverFunctionContext {
	var p = new(RankOverFunctionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_rankOverFunction
	return p
}

func InitEmptyRankOverFunctionContext(p *RankOverFunctionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_rankOverFunction
}

func (*RankOverFunctionContext) IsRankOverFunctionContext() {}

func NewRankOverFunctionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RankOverFunctionContext {
	var p = new(RankOverFunctionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_rankOverFunction

	return p
}

func (s *RankOverFunctionContext) GetParser() antlr.Parser { return s.parser }

func (s *RankOverFunctionContext) RANK_OVER_FUNCTION_NAME() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserRANK_OVER_FUNCTION_NAME, 0)
}

func (s *RankOverFunctionContext) LB() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserLB, 0)
}

func (s *RankOverFunctionContext) AllLSB() []antlr.TerminalNode {
	return s.GetTokens(OutboundAPIParserLSB)
}

func (s *RankOverFunctionContext) LSB(i int) antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserLSB, i)
}

func (s *RankOverFunctionContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(OutboundAPIParserID)
}

func (s *RankOverFunctionContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserID, i)
}

func (s *RankOverFunctionContext) AllRSB() []antlr.TerminalNode {
	return s.GetTokens(OutboundAPIParserRSB)
}

func (s *RankOverFunctionContext) RSB(i int) antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserRSB, i)
}

func (s *RankOverFunctionContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(OutboundAPIParserCOMMA)
}

func (s *RankOverFunctionContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCOMMA, i)
}

func (s *RankOverFunctionContext) AllSORT_ORDER() []antlr.TerminalNode {
	return s.GetTokens(OutboundAPIParserSORT_ORDER)
}

func (s *RankOverFunctionContext) SORT_ORDER(i int) antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserSORT_ORDER, i)
}

func (s *RankOverFunctionContext) RB() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserRB, 0)
}

func (s *RankOverFunctionContext) AllRankOverFilter() []IRankOverFilterContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRankOverFilterContext); ok {
			len++
		}
	}

	tst := make([]IRankOverFilterContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRankOverFilterContext); ok {
			tst[i] = t.(IRankOverFilterContext)
			i++
		}
	}

	return tst
}

func (s *RankOverFunctionContext) RankOverFilter(i int) IRankOverFilterContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRankOverFilterContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRankOverFilterContext)
}

func (s *RankOverFunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RankOverFunctionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RankOverFunctionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterRankOverFunction(s)
	}
}

func (s *RankOverFunctionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitRankOverFunction(s)
	}
}

func (s *RankOverFunctionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitRankOverFunction(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) RankOverFunction() (localctx IRankOverFunctionContext) {
	localctx = NewRankOverFunctionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, OutboundAPIParserRULE_rankOverFunction)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(160)
		p.Match(OutboundAPIParserRANK_OVER_FUNCTION_NAME)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(161)
		p.Match(OutboundAPIParserLB)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(162)
		p.Match(OutboundAPIParserLSB)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(163)
		p.Match(OutboundAPIParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(168)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == OutboundAPIParserCOMMA {
		{
			p.SetState(164)
			p.Match(OutboundAPIParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(165)
			p.Match(OutboundAPIParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(170)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(171)
		p.Match(OutboundAPIParserRSB)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(172)
		p.Match(OutboundAPIParserCOMMA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(173)
		p.Match(OutboundAPIParserLSB)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(174)
		p.Match(OutboundAPIParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(175)
		p.Match(OutboundAPIParserSORT_ORDER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(181)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == OutboundAPIParserCOMMA {
		{
			p.SetState(176)
			p.Match(OutboundAPIParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(177)
			p.Match(OutboundAPIParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(178)
			p.Match(OutboundAPIParserSORT_ORDER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(183)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(184)
		p.Match(OutboundAPIParserRSB)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(189)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == OutboundAPIParserCOMMA {
		{
			p.SetState(185)
			p.Match(OutboundAPIParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(186)
			p.RankOverFilter()
		}

		p.SetState(191)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(192)
		p.Match(OutboundAPIParserRB)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRankOverFilterContext is an interface to support dynamic dispatch.
type IRankOverFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllSIGNED_INTEGER() []antlr.TerminalNode
	SIGNED_INTEGER(i int) antlr.TerminalNode
	LSB() antlr.TerminalNode
	COMMA() antlr.TerminalNode
	RSB() antlr.TerminalNode
	OPEN_FILTER_INTERVAL_MARKER() antlr.TerminalNode

	// IsRankOverFilterContext differentiates from other interfaces.
	IsRankOverFilterContext()
}

type RankOverFilterContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRankOverFilterContext() *RankOverFilterContext {
	var p = new(RankOverFilterContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_rankOverFilter
	return p
}

func InitEmptyRankOverFilterContext(p *RankOverFilterContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_rankOverFilter
}

func (*RankOverFilterContext) IsRankOverFilterContext() {}

func NewRankOverFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RankOverFilterContext {
	var p = new(RankOverFilterContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_rankOverFilter

	return p
}

func (s *RankOverFilterContext) GetParser() antlr.Parser { return s.parser }

func (s *RankOverFilterContext) AllSIGNED_INTEGER() []antlr.TerminalNode {
	return s.GetTokens(OutboundAPIParserSIGNED_INTEGER)
}

func (s *RankOverFilterContext) SIGNED_INTEGER(i int) antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserSIGNED_INTEGER, i)
}

func (s *RankOverFilterContext) LSB() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserLSB, 0)
}

func (s *RankOverFilterContext) COMMA() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCOMMA, 0)
}

func (s *RankOverFilterContext) RSB() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserRSB, 0)
}

func (s *RankOverFilterContext) OPEN_FILTER_INTERVAL_MARKER() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserOPEN_FILTER_INTERVAL_MARKER, 0)
}

func (s *RankOverFilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RankOverFilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RankOverFilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterRankOverFilter(s)
	}
}

func (s *RankOverFilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitRankOverFilter(s)
	}
}

func (s *RankOverFilterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitRankOverFilter(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) RankOverFilter() (localctx IRankOverFilterContext) {
	localctx = NewRankOverFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, OutboundAPIParserRULE_rankOverFilter)
	p.SetState(205)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 18, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(194)
			p.Match(OutboundAPIParserSIGNED_INTEGER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(195)
			p.Match(OutboundAPIParserLSB)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(196)
			p.Match(OutboundAPIParserSIGNED_INTEGER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(197)
			p.Match(OutboundAPIParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(198)
			p.Match(OutboundAPIParserSIGNED_INTEGER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(199)
			p.Match(OutboundAPIParserRSB)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(200)
			p.Match(OutboundAPIParserLSB)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(201)
			p.Match(OutboundAPIParserSIGNED_INTEGER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(202)
			p.Match(OutboundAPIParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(203)
			p.Match(OutboundAPIParserOPEN_FILTER_INTERVAL_MARKER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(204)
			p.Match(OutboundAPIParserRSB)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILatestFunctionContext is an interface to support dynamic dispatch.
type ILatestFunctionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LATEST_FUNCTION_NAME() antlr.TerminalNode
	LB() antlr.TerminalNode
	AllLatestExpression() []ILatestExpressionContext
	LatestExpression(i int) ILatestExpressionContext
	RB() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsLatestFunctionContext differentiates from other interfaces.
	IsLatestFunctionContext()
}

type LatestFunctionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLatestFunctionContext() *LatestFunctionContext {
	var p = new(LatestFunctionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_latestFunction
	return p
}

func InitEmptyLatestFunctionContext(p *LatestFunctionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_latestFunction
}

func (*LatestFunctionContext) IsLatestFunctionContext() {}

func NewLatestFunctionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LatestFunctionContext {
	var p = new(LatestFunctionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_latestFunction

	return p
}

func (s *LatestFunctionContext) GetParser() antlr.Parser { return s.parser }

func (s *LatestFunctionContext) LATEST_FUNCTION_NAME() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserLATEST_FUNCTION_NAME, 0)
}

func (s *LatestFunctionContext) LB() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserLB, 0)
}

func (s *LatestFunctionContext) AllLatestExpression() []ILatestExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILatestExpressionContext); ok {
			len++
		}
	}

	tst := make([]ILatestExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILatestExpressionContext); ok {
			tst[i] = t.(ILatestExpressionContext)
			i++
		}
	}

	return tst
}

func (s *LatestFunctionContext) LatestExpression(i int) ILatestExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILatestExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILatestExpressionContext)
}

func (s *LatestFunctionContext) RB() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserRB, 0)
}

func (s *LatestFunctionContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(OutboundAPIParserCOMMA)
}

func (s *LatestFunctionContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCOMMA, i)
}

func (s *LatestFunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LatestFunctionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LatestFunctionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterLatestFunction(s)
	}
}

func (s *LatestFunctionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitLatestFunction(s)
	}
}

func (s *LatestFunctionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitLatestFunction(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) LatestFunction() (localctx ILatestFunctionContext) {
	localctx = NewLatestFunctionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, OutboundAPIParserRULE_latestFunction)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(207)
		p.Match(OutboundAPIParserLATEST_FUNCTION_NAME)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(208)
		p.Match(OutboundAPIParserLB)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(209)
		p.LatestExpression()
	}
	p.SetState(214)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == OutboundAPIParserCOMMA {
		{
			p.SetState(210)
			p.Match(OutboundAPIParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(211)
			p.LatestExpression()
		}

		p.SetState(216)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(217)
		p.Match(OutboundAPIParserRB)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILatestExpressionContext is an interface to support dynamic dispatch.
type ILatestExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	COMPARISON_OPERATOR() antlr.TerminalNode
	PointInTimeArithmetic() IPointInTimeArithmeticContext
	IN() antlr.TerminalNode
	TimeIntervalArithmetic() ITimeIntervalArithmeticContext
	TimeIntervalToPointInTime() ITimeIntervalToPointInTimeContext
	SIGNED_INTEGER() antlr.TerminalNode
	FLOAT() antlr.TerminalNode

	// IsLatestExpressionContext differentiates from other interfaces.
	IsLatestExpressionContext()
}

type LatestExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLatestExpressionContext() *LatestExpressionContext {
	var p = new(LatestExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_latestExpression
	return p
}

func InitEmptyLatestExpressionContext(p *LatestExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_latestExpression
}

func (*LatestExpressionContext) IsLatestExpressionContext() {}

func NewLatestExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LatestExpressionContext {
	var p = new(LatestExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_latestExpression

	return p
}

func (s *LatestExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *LatestExpressionContext) ID() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserID, 0)
}

func (s *LatestExpressionContext) COMPARISON_OPERATOR() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserCOMPARISON_OPERATOR, 0)
}

func (s *LatestExpressionContext) PointInTimeArithmetic() IPointInTimeArithmeticContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPointInTimeArithmeticContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPointInTimeArithmeticContext)
}

func (s *LatestExpressionContext) IN() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserIN, 0)
}

func (s *LatestExpressionContext) TimeIntervalArithmetic() ITimeIntervalArithmeticContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITimeIntervalArithmeticContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITimeIntervalArithmeticContext)
}

func (s *LatestExpressionContext) TimeIntervalToPointInTime() ITimeIntervalToPointInTimeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITimeIntervalToPointInTimeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITimeIntervalToPointInTimeContext)
}

func (s *LatestExpressionContext) SIGNED_INTEGER() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserSIGNED_INTEGER, 0)
}

func (s *LatestExpressionContext) FLOAT() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserFLOAT, 0)
}

func (s *LatestExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LatestExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LatestExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterLatestExpression(s)
	}
}

func (s *LatestExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitLatestExpression(s)
	}
}

func (s *LatestExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitLatestExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) LatestExpression() (localctx ILatestExpressionContext) {
	localctx = NewLatestExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, OutboundAPIParserRULE_latestExpression)
	var _la int

	p.SetState(231)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 20, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(219)
			p.Match(OutboundAPIParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(220)
			p.Match(OutboundAPIParserCOMPARISON_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(221)
			p.PointInTimeArithmetic()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(222)
			p.Match(OutboundAPIParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(223)
			p.Match(OutboundAPIParserIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(224)
			p.TimeIntervalArithmetic()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(225)
			p.Match(OutboundAPIParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(226)
			p.Match(OutboundAPIParserCOMPARISON_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(227)
			p.TimeIntervalToPointInTime()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(228)
			p.Match(OutboundAPIParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(229)
			p.Match(OutboundAPIParserCOMPARISON_OPERATOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(230)
			_la = p.GetTokenStream().LA(1)

			if !(_la == OutboundAPIParserSIGNED_INTEGER || _la == OutboundAPIParserFLOAT) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IGenericValueContext is an interface to support dynamic dispatch.
type IGenericValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IN() antlr.TerminalNode
	SORT_ORDER() antlr.TerminalNode
	OPEN_FILTER_INTERVAL_MARKER() antlr.TerminalNode
	DATE() antlr.TerminalNode
	TIME() antlr.TerminalNode
	POINT_IN_TIME() antlr.TerminalNode
	TIME_PERIOD() antlr.TerminalNode
	TIME_ZONE_IANA() antlr.TerminalNode
	SIGNED_INTEGER() antlr.TerminalNode
	FLOAT() antlr.TerminalNode
	ID() antlr.TerminalNode
	WORD() antlr.TerminalNode

	// IsGenericValueContext differentiates from other interfaces.
	IsGenericValueContext()
}

type GenericValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGenericValueContext() *GenericValueContext {
	var p = new(GenericValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_genericValue
	return p
}

func InitEmptyGenericValueContext(p *GenericValueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = OutboundAPIParserRULE_genericValue
}

func (*GenericValueContext) IsGenericValueContext() {}

func NewGenericValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GenericValueContext {
	var p = new(GenericValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = OutboundAPIParserRULE_genericValue

	return p
}

func (s *GenericValueContext) GetParser() antlr.Parser { return s.parser }

func (s *GenericValueContext) IN() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserIN, 0)
}

func (s *GenericValueContext) SORT_ORDER() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserSORT_ORDER, 0)
}

func (s *GenericValueContext) OPEN_FILTER_INTERVAL_MARKER() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserOPEN_FILTER_INTERVAL_MARKER, 0)
}

func (s *GenericValueContext) DATE() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserDATE, 0)
}

func (s *GenericValueContext) TIME() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserTIME, 0)
}

func (s *GenericValueContext) POINT_IN_TIME() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserPOINT_IN_TIME, 0)
}

func (s *GenericValueContext) TIME_PERIOD() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserTIME_PERIOD, 0)
}

func (s *GenericValueContext) TIME_ZONE_IANA() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserTIME_ZONE_IANA, 0)
}

func (s *GenericValueContext) SIGNED_INTEGER() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserSIGNED_INTEGER, 0)
}

func (s *GenericValueContext) FLOAT() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserFLOAT, 0)
}

func (s *GenericValueContext) ID() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserID, 0)
}

func (s *GenericValueContext) WORD() antlr.TerminalNode {
	return s.GetToken(OutboundAPIParserWORD, 0)
}

func (s *GenericValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GenericValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GenericValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.EnterGenericValue(s)
	}
}

func (s *GenericValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(OutboundAPIParserListener); ok {
		listenerT.ExitGenericValue(s)
	}
}

func (s *GenericValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case OutboundAPIParserVisitor:
		return t.VisitGenericValue(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *OutboundAPIParser) GenericValue() (localctx IGenericValueContext) {
	localctx = NewGenericValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, OutboundAPIParserRULE_genericValue)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(233)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&8279626612736) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
