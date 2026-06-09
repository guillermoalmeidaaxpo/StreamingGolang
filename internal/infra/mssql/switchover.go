package mssql

import (
	"strings"

	"streaming-golang/internal/domain"
)

type switchoverGroups struct {
	MDSSwitchover  []domain.Mapping
	CMDPSwitchover []domain.Mapping
	NoSwitchover   []domain.Mapping
}

func groupBySwitchover(mappings []domain.Mapping) switchoverGroups {
	var groups switchoverGroups
	for _, mapping := range mappings {
		switch {
		case strings.HasPrefix(strings.ToLower(mapping.SwitchOver), "mds"):
			groups.MDSSwitchover = append(groups.MDSSwitchover, mapping)
		case strings.HasPrefix(strings.ToLower(mapping.SwitchOver), "cmdp"):
			groups.CMDPSwitchover = append(groups.CMDPSwitchover, mapping)
		default:
			groups.NoSwitchover = append(groups.NoSwitchover, mapping)
		}
	}
	return groups
}

func hyperscaleOrOwnID(mapping domain.Mapping) domain.Identifier {
	if mapping.HyperscaleID != nil {
		return *mapping.HyperscaleID
	}
	return mapping.ID
}

func enrichMDSMappings(mdsMappings, originals []domain.Mapping) []domain.Mapping {
	enriched := make([]domain.Mapping, 0, len(mdsMappings))
	for _, mdsMapping := range mdsMappings {
		copy := mdsMapping
		if original, ok := findOriginalForMDSMapping(mdsMapping, originals); ok {
			copy.SwitchOver = original.SwitchOver
			copy.Timezone = firstNonEmpty(copy.Timezone, original.Timezone)
		}
		enriched = append(enriched, copy)
	}
	return enriched
}

func findOriginalForMDSMapping(mdsMapping domain.Mapping, originals []domain.Mapping) (domain.Mapping, bool) {
	for _, original := range originals {
		if original.HyperscaleID != nil && *original.HyperscaleID == mdsMapping.ID {
			return original, true
		}
		if original.HyperscaleID == nil && original.ID == mdsMapping.ID {
			return original, true
		}
	}
	return domain.Mapping{}, false
}

func forceCMDP(mapping domain.Mapping) domain.Mapping {
	mapping.HyperscaleID = nil
	mapping.Source = domain.SourceCMDP
	return mapping
}

func distinctIdentifiers(ids []domain.Identifier) []domain.Identifier {
	seen := make(map[domain.Identifier]struct{}, len(ids))
	result := make([]domain.Identifier, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
