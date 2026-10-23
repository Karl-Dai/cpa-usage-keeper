package poller_test

import (
	"testing"
	"time"

	"cpa-usage-keeper/internal/poller"
)

func newStatusOnlyIngestRunner() *poller.RedisIngestRunner {
	// 只用状态字段做聚合测试，不启动任何后台循环。
	return poller.NewRedisIngestRunner(nil, nil, nil, nil, poller.RedisIngestRunnerConfig{})
}

func TestRedisPollerAggregatesMultipleInstanceIngestStatuses(t *testing.T) {
	earlier := time.Date(2026, 10, 23, 10, 0, 0, 0, time.UTC)
	later := earlier.Add(time.Minute)

	first := newStatusOnlyIngestRunner()
	setUnexportedField(t, first, "running", true)
	setUnexportedField(t, first, "lastRunAt", earlier)
	setUnexportedField(t, first, "lastStatus", "completed")
	setUnexportedField(t, first, "lastWarning", "degraded to http pull")

	second := newStatusOnlyIngestRunner()
	setUnexportedField(t, second, "running", true)
	setUnexportedField(t, second, "lastRunAt", later)
	setUnexportedField(t, second, "lastStatus", "completed_with_warnings")
	setUnexportedField(t, second, "lastError", "instance unreachable")

	aggregated := poller.NewRedisPollerForIngests([]*poller.RedisIngestRunner{first, second}, nil).Status()

	// 汇总状态必须保留每个实例的最新时间、错误和警告，不能只暴露单个实例。
	if !aggregated.Running || !aggregated.SyncRunning {
		t.Fatalf("expected running union across instances, got %+v", aggregated)
	}
	if !aggregated.LastRunAt.Equal(later) || aggregated.LastStatus != "completed_with_warnings" {
		t.Fatalf("expected latest instance status, got %+v", aggregated)
	}
	if aggregated.LastError != "instance unreachable" {
		t.Fatalf("expected instance error to survive aggregation, got %q", aggregated.LastError)
	}
	if aggregated.LastWarning != "degraded to http pull" {
		t.Fatalf("expected instance warning to survive aggregation, got %q", aggregated.LastWarning)
	}
}

func TestRedisPollerAggregatesInstanceErrorsAndProcessStatus(t *testing.T) {
	ingestAt := time.Date(2026, 10, 23, 10, 0, 0, 0, time.UTC)
	processAt := ingestAt.Add(time.Hour)

	first := newStatusOnlyIngestRunner()
	setUnexportedField(t, first, "running", true)
	setUnexportedField(t, first, "lastRunAt", ingestAt)
	setUnexportedField(t, first, "lastStatus", "failed")
	setUnexportedField(t, first, "lastError", "instance a failed")

	second := newStatusOnlyIngestRunner()
	setUnexportedField(t, second, "running", true)
	setUnexportedField(t, second, "lastRunAt", ingestAt)
	setUnexportedField(t, second, "lastStatus", "failed")
	setUnexportedField(t, second, "lastError", "instance b failed")

	process := poller.NewRedisProcessRunner(nil)
	setUnexportedField(t, process, "running", true)
	setUnexportedField(t, process, "processRunning", true)
	setUnexportedField(t, process, "lastRunAt", processAt)
	setUnexportedField(t, process, "lastStatus", "completed")
	setUnexportedField(t, process, "lastError", "process failed")

	aggregated := poller.NewRedisPollerForIngests([]*poller.RedisIngestRunner{first, second}, process).Status()

	if aggregated.LastError != "instance a failed; instance b failed; process failed" {
		t.Fatalf("expected every instance and process error to be joined, got %q", aggregated.LastError)
	}
	if !aggregated.LastRunAt.Equal(processAt) || aggregated.LastStatus != "completed" {
		t.Fatalf("expected process status to win on newer time, got %+v", aggregated)
	}
	// 两个实例都有错误时 ingest 侧 SyncRunning 必须为 false，避免状态接口误报远端可用。
	if ingestOnly := poller.NewRedisPollerForIngests([]*poller.RedisIngestRunner{first, second}, nil).Status(); ingestOnly.SyncRunning {
		t.Fatalf("expected sync running false when every instance failed, got %+v", ingestOnly)
	}
}

func TestRedisPollerSingleInstanceConstructorKeepsExistingBehavior(t *testing.T) {
	ingestAt := time.Date(2026, 10, 23, 10, 0, 0, 0, time.UTC)
	processAt := ingestAt.Add(time.Hour)

	ingest := newStatusOnlyIngestRunner()
	setUnexportedField(t, ingest, "running", true)
	setUnexportedField(t, ingest, "lastRunAt", ingestAt)
	setUnexportedField(t, ingest, "lastStatus", "completed")
	setUnexportedField(t, ingest, "lastError", "ingest failed")

	process := poller.NewRedisProcessRunner(nil)
	setUnexportedField(t, process, "running", true)
	setUnexportedField(t, process, "lastRunAt", processAt)
	setUnexportedField(t, process, "lastStatus", "process-status")

	aggregated := poller.NewRedisPoller(ingest, process).Status()

	if aggregated.LastError != "ingest failed" {
		t.Fatalf("expected single ingest error, got %q", aggregated.LastError)
	}
	if !aggregated.LastRunAt.Equal(processAt) || aggregated.LastStatus != "process-status" {
		t.Fatalf("expected newer process status to win, got %+v", aggregated)
	}

	// nil ingest 仍然返回稳定的空 ingest 状态，保持状态接口兼容。
	if status := poller.NewRedisPoller(nil, nil).Status(); status.Running || status.LastStatus != "" {
		t.Fatalf("expected empty status for nil runners, got %+v", status)
	}
}
