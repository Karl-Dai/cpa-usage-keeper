import { useCallback, useEffect, useLayoutEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import '@/lib/chartjs';
import type { ChartData, ChartOptions, Plugin } from 'chart.js';
import { Scatter } from 'react-chartjs-2';
import type { AnalysisModelSpeedPayload } from '@/lib/types';
import { formatCompactNumber } from '@/utils/usage';
import { getUsageChartTheme, USAGE_CHART_COMPOSITION_COLORS } from '@/utils/usage/chartConfig';
import { formatModelSpeedTime, mergeModelSpeedBuckets, modelSpeedBucketIndex, modelSpeedMean, modelSpeedTooltipRows } from './analysisModelSpeed';
import { ModelSpeedPicker } from './ModelSpeedPicker';
import { ModelSpeedTooltip } from './ModelSpeedTooltip';
import { useModelSpeedSelection } from './modelSpeedSelection';
import panelStyles from '../analysis/AnalysisPanel.module.scss';
import styles from './ModelSpeedChart.module.scss';

type SpeedPoint = { x: number; y: number | null };
type Hover = { source: AnalysisModelSpeedPayload; bucket: number; x: number; y: number; touch: boolean };
const EMPTY_SPEED: AnalysisModelSpeedPayload = { buckets: [], series: [] };
const colorAt = (slot: number | undefined, isDark: boolean) => {
  const color = USAGE_CHART_COMPOSITION_COLORS[slot ?? 0];
  return isDark ? color.light : color.base;
};

export function AnalysisModelSpeedCard({ speed: response, loading, isDark, isMobile, timezone = 'UTC', preferenceScope = 'admin' }: {
  speed?: AnalysisModelSpeedPayload;
  loading: boolean;
  isDark: boolean;
  isMobile: boolean;
  timezone?: string;
  preferenceScope?: string;
}) {
  const { t, i18n } = useTranslation();
  const speed = useMemo(() => mergeModelSpeedBuckets(response ?? EMPTY_SPEED), [response]);
  const { selected, colors, candidates, choose, reset } = useModelSpeedSelection('analysis', preferenceScope, speed.series);
  const theme = useMemo(() => getUsageChartTheme(isDark), [isDark]);
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const touchInput = useRef(false);
  const hoveringTooltip = useRef(false);
  const closeTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const [hover, setHover] = useState<Hover | null>(null);
  const [pickerOpen, setPickerOpen] = useState(false);
  const locale = i18n?.language ?? 'en';
  const models = useMemo(() => new Map(speed.series.map((series) => [series.model, series])), [speed]);
  const labels = useMemo(() => speed.buckets.map((bucket) => formatModelSpeedTime(bucket.start, timezone, locale, true)), [speed, timezone, locale]);
  const chartData = useMemo<ChartData<'scatter', SpeedPoint[]>>(() => ({
    datasets: selected.map((model, modelIndex) => ({
      label: model,
      data: speed.buckets.map((_, bucket) => ({
        x: bucket + (selected.length > 1 ? (modelIndex / (selected.length - 1) - 0.5) * 0.6 : 0),
        y: modelSpeedMean(models.get(model)?.samples[bucket]),
      })),
      backgroundColor: colorAt(colors[model], isDark),
      borderColor: colorAt(colors[model], isDark),
      pointRadius: isMobile ? 2.5 : 4,
      pointHoverRadius: isMobile ? 3.5 : 5,
      borderWidth: 1,
    })),
  }), [selected, speed, models, colors, isDark, isMobile]);
  const cancelClose = useCallback(() => {
    if (closeTimer.current !== null) clearTimeout(closeTimer.current);
    closeTimer.current = null;
  }, []);
  const dismiss = useCallback(() => { cancelClose(); hoveringTooltip.current = false; setHover(null); }, [cancelClose]);
  const scheduleClose = useCallback(() => {
    cancelClose();
    if (touchInput.current || hoveringTooltip.current) return;
    closeTimer.current = setTimeout(dismiss, 160);
  }, [cancelClose, dismiss]);
  useEffect(() => cancelClose, [cancelClose]);
  const changePickerOpen = useCallback((open: boolean) => { setPickerOpen(open); if (open) dismiss(); }, [dismiss]);

  // 命中整个显示时间组，包含点之间的空白与无样本模型；不依赖 scatter 最近单点。
  const handlers = useMemo<Plugin<'scatter'>>(() => ({
    id: 'analysisModelSpeed',
    beforeDatasetsDraw(chart) {
      const zero = chart.scales.y.getPixelForValue(0);
      chart.ctx.save();
      chart.ctx.lineWidth = 1.5;
      chart.data.datasets.forEach((dataset, datasetIndex) => {
        const metadata = chart.getDatasetMeta(datasetIndex);
        chart.ctx.strokeStyle = String(dataset.borderColor);
        metadata.data.forEach((point, index) => {
          if (modelSpeedMean(models.get(selected[datasetIndex])?.samples[index]) === null) return;
          chart.ctx.beginPath(); chart.ctx.moveTo(point.x, zero); chart.ctx.lineTo(point.x, point.y); chart.ctx.stroke();
        });
      });
      chart.ctx.restore();
    },
    afterEvent(chart, { event }) {
      if (pickerOpen || !selected.length || !speed.buckets.length) return;
      if (event.type === 'mouseout') { scheduleClose(); return; }
      if (event.type !== 'mousemove' && event.type !== 'click') return;
      if (touchInput.current && event.type !== 'click') return;
      const { x, y } = event;
      const area = chart.chartArea;
      if (x === null || y === null || x < area.left || x > area.right || y < area.top || y > area.bottom) { scheduleClose(); return; }
      const value = chart.scales.x.getValueForPixel(x);
      if (typeof value !== 'number') return;
      cancelClose();
      canvasRef.current = chart.canvas;
      const bucket = modelSpeedBucketIndex(value, speed.buckets.length);
      setHover((previous) => previous?.source === speed && previous.bucket === bucket && previous.x === x && previous.y === y && previous.touch === touchInput.current
        ? previous : { source: speed, bucket, x, y, touch: touchInput.current });
    },
  }), [models, selected, speed, pickerOpen, cancelClose, scheduleClose]);
  const liveHandlers = useRef(handlers);
  useLayoutEffect(() => { liveHandlers.current = handlers; }, [handlers]);
  // react-chartjs-2 只在建图时注册局部插件，稳定实例从最新 handlers 读取选择和刷新数据。
  const plugin = useMemo<Plugin<'scatter'>>(() => ({
    id: 'analysisModelSpeed',
    beforeDatasetsDraw: (chart, args, options) => liveHandlers.current.beforeDatasetsDraw?.(chart, args, options),
    afterEvent: (chart, args, options) => liveHandlers.current.afterEvent?.(chart, args, options),
  }), []);
  const chartOptions = useMemo<ChartOptions<'scatter'>>(() => ({
    responsive: true,
    maintainAspectRatio: false,
    layout: { padding: { top: 8, right: 8, bottom: 4 } },
    plugins: { legend: { display: false }, tooltip: { enabled: false } },
    scales: {
      x: { type: 'linear', min: -0.5, max: Math.max(0.5, speed.buckets.length - 0.5),
        grid: { color: theme.grid, drawTicks: false }, border: { color: theme.axis },
        afterBuildTicks: (scale) => { scale.ticks = labels.map((_, index) => ({ value: index })); },
        ticks: { color: theme.textSecondary, font: { size: 10 }, maxRotation: 0, autoSkip: false,
          callback: (value) => labels[Number(value)]?.split(/[,，]?\s+/) ?? '' },
      },
      y: { beginAtZero: true, min: 0, grid: { color: theme.grid }, border: { color: theme.axis },
        title: { display: true, text: 't/s', color: theme.textSecondary, font: { size: 10 } },
        ticks: { color: theme.textSecondary, font: { size: 10 }, maxTicksLimit: 5, callback: (value) => formatCompactNumber(Number(value)) },
      },
    },
  }), [speed.buckets.length, labels, theme]);
  const visibleHover = !loading && !pickerOpen && hover?.source === speed && selected.length ? hover : null;
  const hoveredBucket = visibleHover ? speed.buckets[visibleHover.bucket] : null;
  const tooltipRows = visibleHover ? modelSpeedTooltipRows(speed, selected, visibleHover.bucket).map((row) => ({
    model: row.model,
    color: colorAt(colors[row.model], isDark),
    detail: `${row.mean === null ? t('usage_stats.no_data') : `${row.mean.toFixed(1)} t/s`} · ${t('usage_stats.model_speed_samples', { samples: row.count.toLocaleString(locale) })}`,
  })) : [];

  return <section className={`${panelStyles.analysisCard} keeper-card-surface`} data-analysis-model-speed>
    <div className={`${styles.header} keeper-card-heading`}>
      <div><h2 className="keeper-card-title">{t('usage_stats.analysis_model_speed_title')}</h2>
        <p className="keeper-card-subtitle">{t('usage_stats.analysis_model_speed_subtitle')}</p></div>
      <ModelSpeedPicker candidates={candidates} selected={selected} onChange={choose} onReset={reset} disabled={loading} onOpenChange={changePickerOpen} />
    </div>
    <div className={styles.body}>
      {loading ? <div className={`${panelStyles.emptyState} ${styles.empty}`}>{t('common.loading')}</div>
        : selected.length === 0 ? <div className={`${panelStyles.emptyState} ${styles.empty}`}>{t('usage_stats.model_speed_empty_selection')}</div>
        : speed.buckets.length === 0 ? <div className={`${panelStyles.emptyState} ${styles.empty}`}>{t('usage_stats.no_data')}</div>
        : <div className={styles.chart} onPointerDown={(event) => { touchInput.current = event.pointerType === 'touch'; }}
          onPointerMove={(event) => { if (event.pointerType === 'mouse') touchInput.current = false; }}>
          <Scatter data={chartData} options={chartOptions} plugins={[plugin]} />
        </div>}
      {selected.length > 0 && <div className={`${panelStyles.analysisChartLegend} ${styles.legend}`}>
        {selected.map((model) => <span key={model} className={panelStyles.analysisLegendItem} title={model}>
          <span className={panelStyles.analysisLegendDot} style={{ backgroundColor: colorAt(colors[model], isDark) }} aria-hidden="true" />
          <span className={panelStyles.analysisLegendLabel}>{model}</span>
        </span>)}
      </div>}
    </div>
    {visibleHover && hoveredBucket && <ModelSpeedTooltip canvasRef={canvasRef} x={visibleHover.x} y={visibleHover.y} touch={visibleHover.touch}
      title={`${formatModelSpeedTime(hoveredBucket.start, timezone, locale)} – ${formatModelSpeedTime(hoveredBucket.end, timezone, locale)}`}
      rows={tooltipRows} dismiss={dismiss} cancelClose={() => { hoveringTooltip.current = true; cancelClose(); }}
      scheduleClose={() => { hoveringTooltip.current = false; scheduleClose(); }} />}
  </section>;
}
