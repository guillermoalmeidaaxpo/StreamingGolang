lexer grammar OutboundAPILexer;

CI_VALIDITY_PERIOD_START
    : V A L I D I T Y P E R I O D S T A R T
    ;

CI_VALIDITY_PERIOD_END
    : V A L I D I T Y P E R I O D E N D
    ;

CI_INSTANCE_CODE
    : I N S T A N C E C O D E
    ;

CI_BIDID
    : B I D I D
    ;

CI_ISP
    : I S P
    ;

CI_DIRECTION
    : D I R E C T I O N
    ;

CI_STATUS
    : S T A T U S
    ;

CI_OPTION_EXPIRY
    : O P T I O N E X P I R Y
    ;

TIME_INTERVAL_FUNCTION_NAME
    : T I D A Y
    | T I W E E K
    | T I M O N T H
    | T I Q U A R T E R
    | T I Y E A R
    ;

TIME_INTERVAL_GAS_FUNCTION_NAME
    : G A S D A Y E U R O P E
    | G A S W E E K E U R O P E
    | G A S M O N T H E U R O P E
    | G A S Q U A R T E R E U R O P E
    | G A S Y E A R E U R O P E
    | G A S S U M M E R E U R O P E
    | G A S W I N T E R E U R O P E
    ;

TIME_INTERVAL_EXPLICIT_FUNCTION_NAME
    : T I
    ;

LATEST_GLOBAL
    : L A T E S T G L O B A L
    ;

POINT_IN_TIME_FUNCTION_NAME
    : N O W
    ;

POINT_IN_TIME_UTC_FUNCTION_NAME
    : U T C T I M E
    ;

TIME_INTERVAL_TO_POINT_IN_TIME_FUNCTION
    : E N D
    | B E G I N
    ;

RANK_OVER_FUNCTION_NAME
    : R A N K O V E R
    ;

LATEST_FUNCTION_NAME
    : L A T E S T
    ;

IN
    : I N
    ;

COMPARISON_OPERATOR
    : GT
    | GE
    | EQUAL
    | LT
    | LE
    ;

SORT_ORDER
    : D E S C
    | A S C
    ;

OPEN_FILTER_INTERVAL_MARKER
    : L A S T
    ;

EQUAL:              '=' ;
GT:                 '>' ;
LT:                 '<' ;
LE:                 '<=' ;
GE:                 '>=' ;
ADD:                '+' ;
SUB:                '-' ;
MUL:                '*' ;

TIME_ZONE_IANA
    : ( [a-zA-Z]+ '/'* [a-zA-Z_]* '/'+ [a-zA-Z_-]+)
    | C E T | C S T [6] C D T | E E T | E S T | E S T [5] E D T | M E T | M S T | M S T [7] M D T | P S T [8] P D T | H S T
    ;

DATE
    : YEAR '-' MONTH '-' DAY
    ;

TIME
    : HOUR ':' MINUTE ':' SECOND ('.' MS)?
    ;

POINT_IN_TIME
    : DATE 'T' TIME
    ;

TIME_PERIOD
    : 'P' (INTEGER 'Y')? (INTEGER 'M')? (INTEGER 'W')? (INTEGER 'D')? ('T' (INTEGER 'H')? (INTEGER 'M')? (INTEGER 'S')?)?
    ;

fragment YEAR:      [0-9][0-9][0-9][0-9] ;
fragment MONTH:     ( [0][1-9] | [1][0-2] ) ;
fragment DAY:       ( [0][1-9] | [12][0-9] | [3][0-1] ) ;
fragment HOUR:      ( [01][0-9] | [2][0-3] ) ;
fragment MINUTE:    [0-5][0-9] ;
fragment SECOND:    [0-5][0-9] ;
fragment MS:        [0-9][0-9][0-9] ;

fragment A : [aA];
fragment B : [bB];
fragment C : [cC];
fragment D : [dD];
fragment E : [eE];
fragment F : [fF];
fragment G : [gG];
fragment H : [hH];
fragment I : [iI];
fragment J : [jJ];
fragment K : [kK];
fragment L : [lL];
fragment M : [mM];
fragment N : [nN];
fragment O : [oO];
fragment P : [pP];
fragment Q : [qQ];
fragment R : [rR];
fragment S : [sS];
fragment T : [tT];
fragment U : [uU];
fragment V : [vV];
fragment W : [wW];
fragment X : [xX];
fragment Y : [yY];
fragment Z : [zZ];

fragment INTEGER:        [0-9]+ ;

LB:             '(' ;
RB:             ')' ;
LSB:            '[' ;
RSB:            ']' ;
ID:             [a-zA-Z]([a-zA-Z0-9])* ;
SIGNED_INTEGER: [-]? INTEGER ;
FLOAT:          [-]? INTEGER DECIMAL_POINT INTEGER ;
WORD:           [a-zA-Z0-9]([a-zA-Z0-9_.\-+])* ;
QUOTE:          '"' ;
COLON:          ':' ;
COMMA:          ',' ;
SEMICOLON:      ';' ;
DECIMAL_POINT:  '.' ;
WS:             [ \t\r\n]+ -> channel(HIDDEN) ;
ERRORCHAR:      .+? ;
