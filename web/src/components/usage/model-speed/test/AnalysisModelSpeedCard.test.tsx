// @vitest-environment happy-dom
import { act } from 'react';
import { createRoot, type Root } from 'react-dom/client';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import type { Chart, ChartData, ChartOptions, Plugin } from 'chart.js';
import type { AnalysisModelSpeedPayload } from '@/lib/types';
import { modelSpeedStorageKey } from '../modelSpeedSelection';
import { AnalysisModelSpeedCard } from '../AnalysisModelSpeedCard';

let captured: { data: ChartData<'scatter', { x: number; y: number | null }[]>; options: ChartOptions<'scatter'>; plugins: Plugin<'scatter'>[] };
vi.mock('react-chartjs-2', () => ({ Scatter: (props: typeof captured) => { captured = props; return <canvas data-speed-chart />; } }));
vi.mock('@/i18n', () => ({ default: { language: 'en', t: (key: string) => key } }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ i18n: { language: 'en' }, t: (key: string, values?: { samples: string }) => values ? `${key}: ${values.samples}` : key }) }));
globalThis.IS_REACT_ACT_ENVIRONMENT = true;
const speed: AnalysisModelSpeedPayload = {
  buckets: [{ start: '2026-10-08T00:00:00+08:00', end: '2026-10-08T03:00:00+08:00' }, { start: '2026-10-08T03:00:00+08:00', end: '2026-10-08T06:00:00+08:00' }],
  series: [
    { model: 'slow-model-full-name', total_tokens: 1000, speed_sample_count: 3, samples: [{ speed_tps_sum: 120, speed_sample_count: 3, mean_tps: 40 }, { speed_tps_sum: 0, speed_sample_count: 0, mean_tps: null }] },
    { model: 'fast-model-full-name', total_tokens: 100, speed_sample_count: 1, samples: [{ speed_tps_sum: 80, speed_sample_count: 1, mean_tps: 80 }, { speed_tps_sum: 0, speed_sample_count: 0, mean_tps: null }] },
    { model: 'no-speed-model', total_tokens: 2000, speed_sample_count: 0, samples: [{ speed_tps_sum: 0, speed_sample_count: 0, mean_tps: null }, { speed_tps_sum: 0, speed_sample_count: 0, mean_tps: null }] },
  ],
};

describe('Analysis model speed card', () => {
  let root: Root;
  let container: HTMLDivElement;
  beforeEach(() => { localStorage.clear(); container = document.createElement('div'); document.body.append(container); root = createRoot(container); });
  afterEach(async () => { await act(async () => root.unmount()); container.remove(); });
  const render = async (response = speed, granularity: 'hourly' | 'daily' = 'hourly') => { await act(async () => root.render(<AnalysisModelSpeedCard speed={response} granularity={granularity} loading={false} isDark={false} isMobile={false} timezone="Asia/Shanghai" />)); };
  const event = async (x: number, y = 100, type: 'mousemove' | 'click' | 'mouseout' = 'mousemove', plugin = captured.plugins[0]) => {
    const chart = { canvas: container.querySelector('canvas'), chartArea: { left: 0, right: 300, top: 0, bottom: 300 }, scales: { x: { getValueForPixel: (value: number) => value } } } as unknown as Chart<'scatter'>;
    await act(async () => plugin.afterEvent!(chart, { event: { type, x, y, native: null }, replay: false, inChartArea: true }, {}));
  };

  it('plots arithmetic means, null gaps, a shared zero-based axis and stable model colors', async () => {
    await render();
    expect(captured.data.datasets.map((dataset) => dataset.label)).toEqual(['no-speed-model', 'slow-model-full-name', 'fast-model-full-name']);
    expect(captured.data.datasets[0].data[0].y).toBeNull();
    expect(captured.data.datasets[1].data[0].y).toBeCloseTo(40);
    expect(captured.data.datasets[1].data[1].y).toBeNull();
    expect(captured.options.scales?.y?.min).toBe(0);
  });

  it('shows all selected models when hovering a point or empty space in its group and hides beneath the picker', async () => {
    await render();
    await event(0.3);
    const tooltip = document.querySelector<HTMLElement>('[data-model-speed-tooltip]')!;
    expect(tooltip.textContent).toContain('fast-model-full-name');
    expect(tooltip.textContent).toContain('slow-model-full-name');
    expect(tooltip.textContent).toContain('no-speed-model');
    expect(tooltip.textContent).toContain('usage_stats.no_data');
    expect(tooltip.textContent).toContain('80.0 t/s');
    expect(tooltip.textContent).toContain('00:00');
    expect(tooltip.textContent).toContain('03:00');
    expect(Array.from(tooltip.querySelectorAll('li strong:first-child'), (item) => item.textContent)).toEqual(['fast-model-full-name', 'slow-model-full-name', 'no-speed-model']);
    await event(1.25);
    expect(document.querySelector('[data-model-speed-tooltip]')!.querySelectorAll('li')).toHaveLength(3);
    expect(document.querySelector('[data-model-speed-tooltip]')!.textContent).not.toContain('t/s');
    await act(async () => container.querySelector<HTMLButtonElement>('button')!.click());
    expect(document.querySelector('[data-model-speed-picker]')).not.toBeNull();
    expect(document.querySelector('[data-model-speed-tooltip]')).toBeNull();
  });

  it('uses time-only hourly ticks and title, with emphasized speed beside the model and samples below', async () => {
    await render();
    const tick = captured.options.scales?.x?.ticks?.callback as (value: number) => string;
    expect([tick(0), tick(1)]).toEqual(['00:00', '03:00']);
    await event(0.3);
    const tooltip = document.querySelector<HTMLElement>('[data-model-speed-tooltip]')!;
    expect(tooltip.querySelector(':scope > strong')!.textContent).toBe('00:00 – 03:00');
    const value = tooltip.querySelector<HTMLElement>('[data-model-speed-value]')!;
    expect(value.textContent).toBe('80.0 t/s');
    expect(getComputedStyle(value).color).toBe('#111827');
    expect(tooltip.querySelector('li small')!.textContent).toBe('usage_stats.model_speed_samples: 1');
  });

  it('uses short month/day ticks and a date-only tooltip interval for daily response buckets', async () => {
    const daily = structuredClone(speed);
    daily.buckets = [
      { start: '2026-10-08T00:00:00+08:00', end: '2026-10-09T00:00:00+08:00' },
      { start: '2026-10-09T00:00:00+08:00', end: '2026-10-10T00:00:00+08:00' },
    ];
    await render(daily, 'daily');
    const tick = captured.options.scales?.x?.ticks?.callback as (value: number) => string;
    expect([tick(0), tick(1)]).toEqual(['10/8', '10/9']);
    await event(0.3);
    expect(document.querySelector('[data-model-speed-tooltip] > strong')!.textContent).toBe('10/8 – 10/9');
  });

  it('keeps the registered plugin reading current selections and refreshed response data', async () => {
    await render();
    const originalPlugin = captured.plugins[0];
    await act(async () => container.querySelector<HTMLButtonElement>('button')!.click());
    const option = Array.from(document.querySelectorAll<HTMLButtonElement>('[data-model-speed-picker] [role="option"]')).find((item) => item.textContent === 'fast-model-full-name')!;
    await act(async () => option.click());
    await act(async () => document.body.dispatchEvent(new PointerEvent('pointerdown', { bubbles: true })));
    expect(captured.plugins[0]).toBe(originalPlugin);
    await event(0, 100, 'mousemove', originalPlugin);
    expect(document.querySelector('[data-model-speed-tooltip]')!.textContent).not.toContain('fast-model-full-name');
    const refreshed = structuredClone(speed);
    refreshed.series[0].samples[0].speed_tps_sum = 600;
    await render(refreshed);
    expect(captured.plugins[0]).toBe(originalPlugin);
    await event(0, 100, 'mousemove', originalPlugin);
    expect(document.querySelector('[data-model-speed-tooltip]')!.textContent).toContain('200.0 t/s');
  });

  it('lets the pointer pass through the tooltip, follows the next time group and closes on canvas leave', async () => {
    await render();
    await event(0.3);
    const tooltip = document.querySelector<HTMLElement>('[data-model-speed-tooltip]')!;
    expect(getComputedStyle(tooltip).pointerEvents).toBe('none');
    expect(tooltip.querySelector('button')).toBeNull();
    await event(1.25);
    expect(document.querySelector('[data-model-speed-tooltip]')!.textContent).toContain('06:00');
    await event(1.25, 100, 'mouseout');
    expect(document.querySelector('[data-model-speed-tooltip]')).toBeNull();
  });

  it('keeps a touch-selected time group open on canvas leave and dismisses on an outside tap', async () => {
    await render();
    await act(async () => container.querySelector('canvas')!.dispatchEvent(new PointerEvent('pointerdown', { pointerType: 'touch', bubbles: true })));
    await event(0.3, 100, 'click');
    await event(0.3, 100, 'mouseout');
    expect(document.querySelector('[data-model-speed-tooltip]')!.textContent).toContain('80.0 t/s');
    await act(async () => document.body.dispatchEvent(new PointerEvent('pointerdown', { pointerType: 'touch', bubbles: true })));
    expect(document.querySelector('[data-model-speed-tooltip]')).toBeNull();
  });

  it('shows no data when the current range offers no models to choose', async () => {
    await render({ buckets: [], series: [] });
    expect(container.textContent).toContain('usage_stats.no_data');
    expect(container.textContent).not.toContain('usage_stats.model_speed_empty_selection');
  });

  it('preserves explicit empty selection and displays the selection empty state', async () => {
    localStorage.setItem(modelSpeedStorageKey('analysis', 'admin'), JSON.stringify({ version: 1, models: [] }));
    await render();
    expect(container.querySelector('canvas')).toBeNull();
    expect(container.textContent).toContain('usage_stats.model_speed_empty_selection');
  });
});
