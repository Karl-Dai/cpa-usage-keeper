import type { AnalysisModelSpeedPayload, AnalysisModelSpeedSample } from '@/lib/types';

export function modelSpeedMean(sample?: Pick<AnalysisModelSpeedSample, 'speed_tps_sum' | 'speed_sample_count'>): number | null {
  if (!sample || sample.speed_sample_count <= 0) return null;
  const mean = sample.speed_tps_sum / sample.speed_sample_count;
  return Number.isFinite(mean) ? mean : null;
}

// 按服务端真实时间边界合并连续组，始终先加 sum/count；不重新推算自然日长度。
export function mergeModelSpeedBuckets(speed: AnalysisModelSpeedPayload, maximum = 8): AnalysisModelSpeedPayload {
  const count = speed.buckets.length;
  const groups = Math.min(count, maximum);
  if (count <= maximum) return speed;
  const ranges = Array.from({ length: groups }, (_, index) => {
    const start = index * Math.floor(count / groups) + Math.min(index, count % groups);
    return { start, end: start + Math.floor(count / groups) + (index < count % groups ? 1 : 0) };
  });
  return {
    ...speed,
    buckets: ranges.map(({ start, end }) => ({ start: speed.buckets[start].start, end: speed.buckets[end - 1].end })),
    series: speed.series.map((series) => ({
      ...series,
      samples: ranges.map(({ start, end }) => {
        const sample = { speed_tps_sum: 0, speed_sample_count: 0 };
        for (let index = start; index < end; index++) {
          sample.speed_tps_sum += series.samples[index]?.speed_tps_sum ?? 0;
          sample.speed_sample_count += series.samples[index]?.speed_sample_count ?? 0;
        }
        return { ...sample, mean_tps: modelSpeedMean(sample) };
      }),
    })),
  };
}

export function modelSpeedTooltipRows(speed: AnalysisModelSpeedPayload, selected: readonly string[], bucket: number) {
  const models = new Map(speed.series.map((series) => [series.model, series]));
  return selected.map((model) => {
    const sample = models.get(model)?.samples[bucket];
    return { model, mean: modelSpeedMean(sample), count: sample?.speed_sample_count ?? 0 };
  }).sort((left, right) => {
    if (left.mean === null && right.mean !== null) return 1;
    if (right.mean === null && left.mean !== null) return -1;
    return (right.mean ?? 0) - (left.mean ?? 0) || (left.model < right.model ? -1 : left.model > right.model ? 1 : 0);
  });
}

export function modelSpeedBucketIndex(value: number, count: number) {
  return Math.max(0, Math.min(count - 1, Math.round(value)));
}

export function formatModelSpeedTime(timestamp: string, timezone: string, locale: string, short = false, includeSeconds = false) {
  try {
    return new Intl.DateTimeFormat(locale, {
      timeZone: timezone, month: '2-digit', day: '2-digit',
      ...(short ? {} : { year: 'numeric' as const }),
      hour: '2-digit', minute: '2-digit', hourCycle: 'h23',
      ...(includeSeconds ? { second: '2-digit' as const } : {}),
    }).format(new Date(timestamp));
  } catch { return timestamp.replace('T', ' ').slice(0, 16); }
}
