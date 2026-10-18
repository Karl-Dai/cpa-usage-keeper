package test

import (
	"testing"
	"time"

	"cpa-usage-keeper/internal/entities"
	"cpa-usage-keeper/internal/repository"
	"cpa-usage-keeper/internal/repository/dto"
	"cpa-usage-keeper/internal/timeutil"
)

func TestAnalysisCalendarDayKeepsHourlyDataAndEndingSentinel(t *testing.T) {
	for _, tc := range []struct {
		zone       string
		month      time.Month
		day, hours int
	}{
		{"Asia/Shanghai", 3, 8, 24},
		{"America/Los_Angeles", 3, 8, 23},
		{"America/Los_Angeles", 11, 1, 25},
		{"Europe/Berlin", 10, 25, 25},
		{"America/Santiago", 9, 6, 23},
	} {
		t.Run(tc.zone+"/"+tc.month.String(), func(t *testing.T) {
			location := setAnalysisLatencyTestTimezone(t, tc.zone)
			anchor := time.Date(2026, tc.month, tc.day, 12, 0, 0, 0, location)
			rangeValue, err := timeutil.ParseUsageQueryRange("today", "", "", "", anchor)
			if err != nil {
				t.Fatal(err)
			}
			start, end := rangeValue.StartTime, rangeValue.EndTime
			db := openTestDatabase(t)
			if err := db.Create(&entities.CPAAPIKey{APIKey: "calendar-key"}).Error; err != nil {
				t.Fatal(err)
			}
			cost, unavailable := 0.0, int64(0)
			var hourly []entities.UsageOverviewHourlyStat
			var events []entities.UsageEvent
			for i := 0; i < tc.hours; i++ {
				at := start.Add(time.Duration(i) * time.Hour)
				hourly = append(hourly, entities.UsageOverviewHourlyStat{BucketStart: at, APIGroupKey: "calendar-key", Model: "model", RequestCount: 1, TotalTokens: 100, CostUSD: &cost, UnavailableCostCount: &unavailable})
				events = append(events, validLatencyEvent(int64(i+1), "calendar-key", at.Add(time.Minute), int64(i+1), int64(i+1)*10))
			}
			if err := db.Create(&hourly).Error; err != nil {
				t.Fatal(err)
			}
			seedAnalysisLatencyRows(t, db, end, entities.UsageLatencyBucketHour, events)
			filter := dto.UsageQueryFilter{Range: "today", StartTime: &start, EndTime: &end}
			analysis, err := repository.BuildAnalysisWithFilter(db, filter)
			if err != nil {
				t.Fatal(err)
			}
			if analysis.Granularity != dto.AnalysisGranularityHourly || len(analysis.TokenUsage) != tc.hours+1 {
				t.Fatalf("want %d hourly rows including sentinel, got %s/%d", tc.hours+1, analysis.Granularity, len(analysis.TokenUsage))
			}
			for i, point := range analysis.TokenUsage {
				want := int64(100)
				if i == tc.hours {
					want = 0
				}
				if point.TotalTokens != want || !point.Bucket.Equal(start.Add(time.Duration(i)*time.Hour)) {
					t.Fatalf("wrong point[%d]: %+v", i, point)
				}
			}
			diagnostics, err := repository.BuildAnalysisLatencyDiagnosticsWithFilter(db, filter)
			if err != nil {
				t.Fatal(err)
			}
			if diagnostics.TotalPoints != int64(tc.hours) {
				t.Fatalf("lost hourly diagnostics: got %d, want %d", diagnostics.TotalPoints, tc.hours)
			}
		})
	}
}

func TestAnalysisLatencyReadsExistingCalendarAndLocalHourBuckets(t *testing.T) {
	for _, tc := range []struct{ zone, start, end, unit string }{
		{"America/Santiago", "2026-09-06T01:00:00-03:00", "2026-09-07T00:00:00-03:00", "day"},
		{"Europe/Berlin", "2026-10-25T02:00:00+02:00", "2026-10-25T02:00:00+01:00", "hour"},
		{"Asia/Kathmandu", "2026-10-25T02:00:00+05:45", "2026-10-25T03:00:00+05:45", "hour"},
	} {
		t.Run(tc.zone, func(t *testing.T) {
			location := setAnalysisLatencyTestTimezone(t, tc.zone)
			start, err := time.Parse(time.RFC3339, tc.start)
			if err != nil {
				t.Fatal(err)
			}
			end, err := time.Parse(time.RFC3339, tc.end)
			if err != nil {
				t.Fatal(err)
			}
			start, end = start.In(location), end.In(location)
			db := openTestDatabase(t)
			kind := entities.UsageLatencyBucketHour
			if tc.unit == "day" {
				kind = entities.UsageLatencyBucketDay
			}
			seedAnalysisLatencyRows(t, db, end.Add(time.Hour), kind, []entities.UsageEvent{
				validLatencyEvent(1, "key", start.Add(time.Minute), 100, 1000),
				validLatencyEvent(2, "key", end.Add(time.Minute), 200, 2000),
			})
			diagnostics, err := repository.BuildAnalysisLatencyDiagnosticsWithFilter(db, dto.UsageQueryFilter{Range: "custom", CustomUnit: tc.unit, StartTime: &start, EndTime: &end, EndExclusive: true})
			if err != nil {
				t.Fatal(err)
			}
			assertSingleAnalysisLatencyPoint(t, diagnostics, 100, 1000)
		})
	}
}
