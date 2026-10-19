package test

import (
	"fmt"
	"testing"
	"time"

	"cpa-usage-keeper/internal/entities"
	"cpa-usage-keeper/internal/repository"
	"cpa-usage-keeper/internal/repository/dto"
)

func TestOverviewRawBoundariesSeparateRepeatedClockTimes(t *testing.T) {
	for _, tc := range []struct {
		name, zone, start, end string
		events                 []string
		want                   int64
	}{
		{"first repeated hour", "Europe/Berlin", "2026-10-25T01:30:00+02:00", "2026-10-25T02:30:00+02:00", []string{"2026-10-25T02:10:00+02:00", "2026-10-25T02:10:00+01:00"}, 1},
		{"crossing backward clock", "Europe/Berlin", "2026-10-25T02:30:00+02:00", "2026-10-25T02:30:00+01:00", []string{"2026-10-25T02:10:00+02:00", "2026-10-25T02:40:00+02:00", "2026-10-25T02:10:00+01:00", "2026-10-25T02:40:00+01:00"}, 2},
		{"Shanghai subsecond endpoints", "Asia/Shanghai", "2026-10-25T02:30:00.5+08:00", "2026-10-25T03:30:00.5+08:00", []string{"2026-10-25T02:30:00.4+08:00", "2026-10-25T02:30:00.5+08:00", "2026-10-25T03:30:00.4+08:00", "2026-10-25T03:30:00.5+08:00"}, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			location := setAnalysisLatencyTestTimezone(t, tc.zone)
			parse := func(text string) time.Time {
				at, err := time.Parse(time.RFC3339Nano, text)
				if err != nil {
					t.Fatal(err)
				}
				return at.In(location)
			}
			start, end := parse(tc.start), parse(tc.end)
			now := end.Add(48 * time.Hour)
			db := openTestDatabase(t)
			cost, available := 0.5, true
			for i, at := range tc.events {
				row := entities.UsageEvent{EventKey: fmt.Sprintf("boundary-%d", i), APIGroupKey: "key", Model: "model", Timestamp: parse(at), TotalTokens: 10, CostUSD: &cost, CostAvailable: &available}
				if err := db.Create(&row).Error; err != nil {
					t.Fatal(err)
				}
			}
			filter := dto.UsageQueryFilter{Range: "1h", StartTime: &start, EndTime: &end, EndExclusive: true, QueryNow: &now}
			result, err := repository.BuildUsageOverviewWithFilterAndRecentCache(db, filter, nil)
			if err != nil {
				t.Fatal(err)
			}
			if result.Usage.TotalRequests != tc.want || result.Usage.TotalTokens != tc.want*10 {
				t.Fatalf("boundary leaked or lost events: %+v", result.Usage)
			}
			assertFloatClose(t, result.Summary.TotalCost, float64(tc.want)*0.5)
			comparison, err := repository.BuildUsageOverviewComparisonsWithFilterAndRecentCache(db, filter, nil)
			if err != nil {
				t.Fatal(err)
			}
			model := comparison.Comparisons.Models["model"]
			if model == nil || model.Requests != tc.want || model.TotalTokens != tc.want*10 {
				t.Fatalf("wrong comparison boundary: %+v", comparison.Comparisons)
			}
		})
	}
}
