package mssql

import (
	"database/sql"
	"testing"

	"streaming-golang/internal/domain"
)

func TestBuildDomainMappingsGroupsColumnsByIdentifier(t *testing.T) {
	rows := []mappingRow{
		{
			TimeSeriesID:    10,
			CMDPViewName:    sql.NullString{String: "[ACCESS].[Data_EodSettlement]", Valid: true},
			MDSDataCategory: "Curve",
			Resolution:      sql.NullString{String: "PT1H", Valid: true},
			CMDPColumnName:  "QuoteTime",
			MDSColumnName:   "ReferenceTime",
			DataType:        "datetimeoffset",
			IsProjectable:   sql.NullBool{Bool: true, Valid: true},
			IsKey:           sql.NullBool{Bool: true, Valid: true},
			IndexField:      sql.NullString{String: "QuoteDateIndex_FID", Valid: true},
			SplitQuery:      sql.NullBool{Bool: true, Valid: true},
		},
		{
			TimeSeriesID:    10,
			CMDPViewName:    sql.NullString{String: "[ACCESS].[Data_EodSettlement]", Valid: true},
			MDSDataCategory: "Curve",
			CMDPColumnName:  "settlement price",
			MDSColumnName:   "SettlementPrice",
			DataType:        "real",
			IsProjectable:   sql.NullBool{Bool: true, Valid: true},
		},
	}

	mappings := buildDomainMappings(rows, domain.Curves)
	if len(mappings) != 1 {
		t.Fatalf("expected one mapping, got %d", len(mappings))
	}
	if mappings[0].ID != 10 {
		t.Fatalf("unexpected id: %d", mappings[0].ID)
	}
	if mappings[0].DataCategory != domain.Curves {
		t.Fatalf("unexpected category: %q", mappings[0].DataCategory)
	}
	if mappings[0].Source != domain.SourceCMDP {
		t.Fatalf("unexpected source: %q", mappings[0].Source)
	}
	if len(mappings[0].Columns) != 2 {
		t.Fatalf("expected two columns, got %d", len(mappings[0].Columns))
	}
}

func TestBuildDomainMappingsDetectsCassandraSource(t *testing.T) {
	rows := []mappingRow{
		{
			TimeSeriesID:    20,
			MDSDataCategory: "TimeSeries",
			CMDPColumnName:  "QuoteTime",
			MDSColumnName:   "ReferenceTime",
			DataType:        "datetimeoffset",
			CassandraID:     sql.NullString{String: "cass-id", Valid: true},
		},
	}

	mappings := buildDomainMappings(rows, domain.TimeSeries)
	if mappings[0].Source != domain.SourceCassandra {
		t.Fatalf("expected cassandra source, got %q", mappings[0].Source)
	}
}

func TestGroupBySwitchover(t *testing.T) {
	groups := groupBySwitchover([]domain.Mapping{
		{ID: 1, SwitchOver: "MDS:2024-01-01"},
		{ID: 2, SwitchOver: "CMDP:2024-01-01"},
		{ID: 3},
	})

	if len(groups.MDSSwitchover) != 1 {
		t.Fatalf("expected one MDS switchover mapping, got %d", len(groups.MDSSwitchover))
	}
	if len(groups.CMDPSwitchover) != 1 {
		t.Fatalf("expected one CMDP switchover mapping, got %d", len(groups.CMDPSwitchover))
	}
	if len(groups.NoSwitchover) != 1 {
		t.Fatalf("expected one non-switchover mapping, got %d", len(groups.NoSwitchover))
	}
}

func TestEnrichMDSMappingsCopiesSwitchoverFromOriginalHyperscaleMapping(t *testing.T) {
	hyperscaleID := domain.Identifier(200)
	originals := []domain.Mapping{{
		ID:           100,
		HyperscaleID: &hyperscaleID,
		SwitchOver:   "MDS:2024-01-01",
		Timezone:     "CET",
	}}
	mdsMappings := []domain.Mapping{{
		ID:     200,
		Source: domain.SourceHyperscale,
	}}

	enriched := enrichMDSMappings(mdsMappings, originals)
	if enriched[0].SwitchOver != "MDS:2024-01-01" {
		t.Fatalf("expected switchover to be copied, got %q", enriched[0].SwitchOver)
	}
	if enriched[0].Timezone != "CET" {
		t.Fatalf("expected timezone to be copied, got %q", enriched[0].Timezone)
	}
}

func TestForceCMDPClearsHyperscaleID(t *testing.T) {
	hyperscaleID := domain.Identifier(200)
	mapping := forceCMDP(domain.Mapping{
		ID:           100,
		Source:       domain.SourceHyperscale,
		HyperscaleID: &hyperscaleID,
	})

	if mapping.HyperscaleID != nil {
		t.Fatal("expected hyperscale id to be cleared")
	}
	if mapping.Source != domain.SourceCMDP {
		t.Fatalf("expected source CMDP, got %q", mapping.Source)
	}
}
