package mssql

import (
	"context"
	"strings"
	"testing"
	"time"

	"streaming-golang/internal/domain"
)

func TestCMDPQueryBuilderBuildsStatementFromMappingsFiltersAndSplitRange(t *testing.T) {
	queryBuilder := NewCMDPQueryBuilder()
	keyOrder := 1
	valueOrder := 1

	queries, err := queryBuilder.BuildQueries(context.Background(), domain.Command{
		DataCategory: domain.Curves,
		Filters: domain.FilterSet{Nodes: []domain.FilterNode{
			domain.ComparisonFilter{
				Field:    "ReferenceTime",
				Operator: "in",
				Value: domain.FilterValue{
					Kind:  domain.FilterValueTimeInterval,
					Raw:   "ti(2023-01-01T00:00:00,2023-01-03T00:00:00)",
					Start: "2023-01-01T00:00:00",
					End:   "2023-01-03T00:00:00",
				},
			},
			domain.ComparisonFilter{
				Field:    "SettlementPrice",
				Operator: ">=",
				Value: domain.FilterValue{
					Kind: domain.FilterValueNumber,
					Raw:  "42.5",
				},
			},
		}},
		Mappings: []domain.Mapping{{
			ID:           536013751,
			DataCategory: domain.Curves,
			Source:       domain.SourceCMDP,
			ViewName:     "ACCESS.Data_PriceModelled",
			IndexField:   "QuoteDateIndex_FID",
			Columns: []domain.ColumnMapping{
				{MDSName: "ReferenceTime", SourceName: "QuoteTime", IsKey: true, IsProjectable: true, KeyColumnOrdering: &keyOrder},
				{MDSName: "SettlementPrice", SourceName: "settlement price", IsProjectable: true, ValueColumnOrdering: &valueOrder},
			},
		}},
		IndexRange: &domain.IndexRange{Start: 20221231, End: 20230104},
	})
	if err != nil {
		t.Fatalf("build queries failed: %v", err)
	}
	if len(queries) != 1 {
		t.Fatalf("query count = %d, want 1", len(queries))
	}

	statement := queries[0].Statement
	assertContains(t, statement, "SELECT [d].[QuoteTime] AS [ReferenceTime], [d].[settlement price] AS [SettlementPrice]")
	assertContains(t, statement, "FROM [ACCESS].[Data_PriceModelled] AS [d]")
	assertContains(t, statement, "[d].[TimeSeries_FID] = @id")
	assertContains(t, statement, "([d].[QuoteTime] >= @p0 AND [d].[QuoteTime] <= @p1)")
	assertContains(t, statement, "[d].[settlement price] >= @p2")
	assertContains(t, statement, "[d].[QuoteDateIndex_FID] >= @indexStart")
	assertContains(t, statement, "[d].[QuoteDateIndex_FID] <= @indexEnd")
	assertContains(t, statement, "ORDER BY [d].[QuoteTime]")

	if queries[0].Parameters["id"] != int64(536013751) {
		t.Fatalf("id parameter = %#v", queries[0].Parameters["id"])
	}
	if queries[0].Parameters["indexStart"] != 20221231 || queries[0].Parameters["indexEnd"] != 20230104 {
		t.Fatalf("index parameters = %#v", queries[0].Parameters)
	}
	if got, want := queries[0].Parameters["p0"], time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC); got != want {
		t.Fatalf("p0 = %#v, want %#v", got, want)
	}
	if got, want := queries[0].Parameters["p2"], 42.5; got != want {
		t.Fatalf("p2 = %#v, want %#v", got, want)
	}
}

func TestCMDPQueryBuilderRejectsUnsupportedLatestFilters(t *testing.T) {
	queryBuilder := NewCMDPQueryBuilder()

	_, err := queryBuilder.BuildQueries(context.Background(), domain.Command{
		DataCategory: domain.Curves,
		Filters: domain.FilterSet{Nodes: []domain.FilterNode{
			domain.ComparisonFilter{
				Field:    "ReferenceTime",
				Operator: "=",
				Value: domain.FilterValue{
					Kind: domain.FilterValueLatestGlobal,
					Raw:  "latestGlobal()",
				},
			},
		}},
		Mappings: []domain.Mapping{{
			ID:       1,
			Source:   domain.SourceCMDP,
			ViewName: "ACCESS.Data_PriceModelled",
			Columns: []domain.ColumnMapping{
				{MDSName: "ReferenceTime", SourceName: "QuoteTime", IsKey: true},
			},
		}},
	})
	if err == nil {
		t.Fatal("expected latest filter to fail")
	}
}

func TestHyperscaleQueryBuilderBuildsStatementFromMDSMappings(t *testing.T) {
	queryBuilder := NewHyperscaleQueryBuilder()
	keyOrder := 1
	valueOrder := 1

	queries, err := queryBuilder.BuildQueries(context.Background(), domain.Command{
		DataCategory: domain.Curves,
		Filters: domain.FilterSet{Nodes: []domain.FilterNode{
			domain.ComparisonFilter{
				Field:    "ReferenceTime",
				Operator: ">=",
				Value: domain.FilterValue{
					Kind: domain.FilterValuePointInTime,
					Raw:  "2025-08-23T00:00:00",
				},
			},
		}},
		Mappings: []domain.Mapping{{
			ID:           488109751,
			DataCategory: domain.Curves,
			Source:       domain.SourceHyperscale,
			Columns: []domain.ColumnMapping{
				{MDSName: "Identifier", SourceName: "Identifier", IsKey: true, IsProjectable: false, KeyColumnOrdering: &keyOrder},
				{MDSName: "ReferenceTime", SourceName: "ReferenceTime", IsKey: true, IsProjectable: false},
				{MDSName: "Value", SourceName: "Value", DataType: "number", IsProjectable: true, ValueColumnOrdering: &valueOrder},
			},
		}},
	})
	if err != nil {
		t.Fatalf("build queries failed: %v", err)
	}
	if len(queries) != 1 {
		t.Fatalf("query count = %d, want 1", len(queries))
	}

	statement := queries[0].Statement
	assertContains(t, statement, "FROM [Api].[VI_CurveLatestVersion] AS [d]")
	assertContains(t, statement, "[d].[MdoId] = @id")
	assertContains(t, statement, "[d].[ReferenceTime] >= @p0")
	assertContains(t, statement, "CAST(JSON_VALUE([d].[CurveValue], '$.\"Value\"') AS FLOAT) AS [Value]")
	assertContains(t, statement, "[d].[Deleted] = 0")
	assertNotContains(t, statement, "[d].[Identifier]")

	if queries[0].Source != domain.SourceHyperscale {
		t.Fatalf("source = %q, want hyperscale", queries[0].Source)
	}
	if queries[0].Parameters["id"] != int64(488109751) {
		t.Fatalf("id parameter = %#v", queries[0].Parameters["id"])
	}
}

func TestHyperscaleQueryBuilderIncludesIdentifierForCSV(t *testing.T) {
	queryBuilder := NewHyperscaleQueryBuilder()
	keyOrder := 1

	queries, err := queryBuilder.BuildQueries(context.Background(), domain.Command{
		DataCategory:      domain.Surfaces,
		IncludeIdentifier: true,
		Mappings: []domain.Mapping{{
			ID:           1000000001,
			DataCategory: domain.Surfaces,
			Source:       domain.SourceHyperscale,
			Columns: []domain.ColumnMapping{
				{MDSName: "Identifier", SourceName: "Identifier", IsKey: true, KeyColumnOrdering: &keyOrder},
				{MDSName: "ReferenceTime", SourceName: "ReferenceTime", IsKey: true},
			},
		}},
	})
	if err != nil {
		t.Fatalf("build queries failed: %v", err)
	}
	if len(queries) != 1 {
		t.Fatalf("query count = %d, want 1", len(queries))
	}

	statement := queries[0].Statement
	assertContains(t, statement, "SELECT [d].[MdoId] AS [Identifier], [d].[ReferenceTime]")
	assertContains(t, statement, "FROM [Api].[VI_SurfaceLatestVersion] AS [d]")
	assertContains(t, statement, "ORDER BY [d].[MdoId]")
}

func TestHyperscaleQueryBuilderBuildsLatestGlobalTimeseriesStatementLikeCSharp(t *testing.T) {
	queryBuilder := NewHyperscaleQueryBuilder()
	keyOrder := 1
	valueOrder := 1

	queries, err := queryBuilder.BuildQueries(context.Background(), domain.Command{
		DataCategory:        domain.TimeSeries,
		Columns:             []string{"CreatedOn"},
		LatestReferenceTime: true,
		IncludeIdentifier:   false,
		IncludeDeleted:      false,
		Mappings: []domain.Mapping{{
			ID:           504078501,
			DataCategory: domain.TimeSeries,
			Source:       domain.SourceHyperscale,
			Views: domain.MappingViews{
				LatestReferenceTime:              "Api.VI_TimeseriesLatestVersionLatestReferenceTime",
				LatestReferenceTimeWithCreatedOn: "Api.VI_TimeseriesLatestVersionLatestReferenceTimeWithCreatedOn",
			},
			Columns: []domain.ColumnMapping{
				{MDSName: "Identifier", SourceName: "Identifier", IsKey: true, KeyColumnOrdering: &keyOrder},
				{MDSName: "ReferenceTime", SourceName: "ReferenceTime", IsKey: true, KeyColumnOrdering: &keyOrder},
				{MDSName: "Value", SourceName: "Value", DataType: "number", IsProjectable: true, ValueColumnOrdering: &valueOrder},
			},
		}},
	})
	if err != nil {
		t.Fatalf("build queries failed: %v", err)
	}
	if len(queries) != 1 {
		t.Fatalf("query count = %d, want 1", len(queries))
	}

	statement := queries[0].Statement
	assertContains(t, statement, `SELECT [d].[ReferenceTime], CAST(JSON_VALUE([d].[TimeSeriesValue], '$."Value"') AS FLOAT) AS [Property0], [d].[CreatedOn]`)
	assertContains(t, statement, "FROM [Api].[VI_TimeseriesLatestVersionLatestReferenceTime] AS [d]")
	assertContains(t, statement, "[d].[MdoId] = @id")
	assertContains(t, statement, "[d].[Deleted] = 0")
	assertContains(t, statement, "ORDER BY [d].[ReferenceTime]")
	assertNotContains(t, statement, "[d].[MdoId] AS [Identifier]")

	if queries[0].Parameters["id"] != int64(504078501) {
		t.Fatalf("id parameter = %#v", queries[0].Parameters["id"])
	}
}

func TestHyperscaleQueryBuilderUsesCreatedOnLatestReferenceTimeTVF(t *testing.T) {
	queryBuilder := NewHyperscaleQueryBuilder()
	versionAsOf := time.Date(2025, 5, 5, 7, 0, 0, 0, time.UTC)

	queries, err := queryBuilder.BuildQueries(context.Background(), domain.Command{
		DataCategory:        domain.TimeSeries,
		Columns:             []string{"ReferenceTime", "Value", "CreatedOn"},
		VersionAsOf:         &versionAsOf,
		LatestReferenceTime: true,
		Mappings: []domain.Mapping{{
			ID:           504078501,
			DataCategory: domain.TimeSeries,
			Source:       domain.SourceHyperscale,
			Views: domain.MappingViews{
				GetByCreatedOn:                    "Api.TVF_GetTimeseriesByCreatedOn",
				GetByCreatedOnLatestReferenceTime: "Api.TVF_GetTimeseriesByCreatedOnLatestReferenceTime",
			},
			Columns: []domain.ColumnMapping{
				{MDSName: "ReferenceTime", SourceName: "ReferenceTime", IsKey: true},
				{MDSName: "Value", SourceName: "Value", DataType: "number", IsProjectable: true},
			},
		}},
	})
	if err != nil {
		t.Fatalf("build queries failed: %v", err)
	}
	if len(queries) != 1 {
		t.Fatalf("query count = %d, want 1", len(queries))
	}

	statement := queries[0].Statement
	assertContains(t, statement, "FROM [Api].[TVF_GetTimeseriesByCreatedOnLatestReferenceTime](@MdoId, @CreatedOn, @IncludeDeleted) AS [d]")
	assertContains(t, statement, `CAST(JSON_VALUE([d].[TimeSeriesValue], '$."Value"') AS FLOAT) AS [Value]`)
	assertContains(t, statement, "[d].[CreatedOn]")
	if queries[0].Parameters["MdoId"] != int64(504078501) {
		t.Fatalf("MdoId parameter = %#v", queries[0].Parameters["MdoId"])
	}
}

func assertContains(t *testing.T, text, substring string) {
	t.Helper()
	if !strings.Contains(text, substring) {
		t.Fatalf("expected statement to contain %q:\n%s", substring, text)
	}
}

func assertNotContains(t *testing.T, text, substring string) {
	t.Helper()
	if strings.Contains(text, substring) {
		t.Fatalf("expected statement not to contain %q:\n%s", substring, text)
	}
}
