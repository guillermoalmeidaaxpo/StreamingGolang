package httpapi

var transactionalRoutes = []string{
	"/curves", "/timeseries", "/surfaces",
	"/design/curves", "/design/timeseries", "/design/surfaces",
	"/validation/curves", "/validation/timeseries", "/validation/surfaces",
	"/productive/curves", "/productive/timeseries", "/productive/surfaces",
	"/migration/curves", "/migration/timeseries", "/migration/surfaces",
}

var transactionalStreamingRoutes = []string{
	"/curves/streaming", "/surfaces/streaming", "/timeseries/streaming",
	"/design/curves/streaming", "/design/surfaces/streaming", "/design/timeseries/streaming",
	"/validation/curves/streaming", "/validation/surfaces/streaming", "/validation/timeseries/streaming",
	"/productive/curves/streaming", "/productive/surfaces/streaming", "/productive/timeseries/streaming",
	"/migration/curves/streaming", "/migration/timeseries/streaming", "/migration/surfaces/streaming",
}

var genericRoutes = []string{
	"/generic", "/design/generic", "/validation/generic", "/productive/generic", "/migration/generic",
}

var genericStreamingRoutes = []string{
	"/generic/streaming", "/design/generic/streaming", "/validation/generic/streaming", "/productive/generic/streaming", "/migration/generic/streaming",
}

var liteRoutes = []string{
	"/lite", "/design/lite", "/validation/lite", "/productive/lite",
}

var dataTraceRoutes = []string{
	"/datatrace", "/design/datatrace", "/validation/datatrace", "/productive/datatrace",
}

var metadataRoutes = []string{
	"/curves/metadata", "/timeseries/metadata", "/surfaces/metadata",
	"/design/curves/metadata", "/design/timeseries/metadata", "/design/surfaces/metadata",
	"/validation/curves/metadata", "/validation/timeseries/metadata", "/validation/surfaces/metadata",
	"/productive/curves/metadata", "/productive/timeseries/metadata", "/productive/surfaces/metadata",
}

var metadataRangeRoutes = []string{
	"/curves/metadata/range", "/timeseries/metadata/range", "/surfaces/metadata/range",
	"/design/curves/metadata/range", "/design/timeseries/metadata/range", "/design/surfaces/metadata/range",
	"/validation/curves/metadata/range", "/validation/timeseries/metadata/range", "/validation/surfaces/metadata/range",
	"/productive/curves/metadata/range", "/productive/timeseries/metadata/range", "/productive/surfaces/metadata/range",
}

var mesapGenericRoutes = []string{
	"/mesaptransition/generic", "/validation/mesaptransition/generic", "/productive/mesaptransition/generic",
}
