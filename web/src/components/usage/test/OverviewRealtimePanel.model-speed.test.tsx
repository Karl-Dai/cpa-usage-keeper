// @vitest-environment happy-dom
import { act, type ComponentProps } from 'react';
import { createRoot } from 'react-dom/client';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import type { OverviewRealtimeBlock } from '@/lib/types';
import { OverviewRealtimePanel } from '../OverviewRealtimePanel';

let speedProps: ComponentProps<typeof import('../model-speed/RealtimeModelSpeedCard').RealtimeModelSpeedCard>;
vi.mock('../model-speed/RealtimeModelSpeedCard', () => ({ RealtimeModelSpeedCard: (props: typeof speedProps) => { speedProps = props; return <section data-model-speed-card>Speed</section>; } }));
vi.mock('react-chartjs-2', () => ({ Line: () => <canvas />, Chart: () => <canvas />, Bar: () => <canvas />, Scatter: () => <canvas /> }));
vi.mock('@/i18n', () => ({ default: { language: 'en', t: (key: string) => key } }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
globalThis.IS_REACT_ACT_ENVIRONMENT = true;
let root: ReturnType<typeof createRoot>;
let container: HTMLDivElement;
beforeEach(() => { container = document.createElement('div'); document.body.append(container); root = createRoot(container); });
afterEach(async () => { await act(async () => root.unmount()); container.remove(); });
const response: OverviewRealtimeBlock = { window: '15m', bucket_seconds: 30, token_velocity: [], request_level: [], cache_level: [], current_usage: { models: [], api_keys: [], auth_files: [], ai_providers: [] }, model_speed: { available: false, series: [] }, timezone: 'Asia/Shanghai' };

it('pairs model speed with Cache Level and forwards the existing response and viewer scope', async () => {
  await act(async () => root.render(<OverviewRealtimePanel realtime={response} loading={false} window="15m" onWindowChange={vi.fn()} isDark isMobile={false} modelSpeedPreferenceScope="key-viewer:42" modelSpeedQueryKey="15m:key42" />));
  const speed = container.querySelector('[data-model-speed-card]')!;
  expect(speed.nextElementSibling?.textContent).toContain('usage_stats.overview_realtime_cache_level');
  expect(speed.parentElement?.children).toHaveLength(2);
  expect(speedProps.speed).toBe(response.model_speed);
  expect(speedProps.timezone).toBe('Asia/Shanghai');
  expect(speedProps.preferenceScope).toBe('key-viewer:42');
  expect(speedProps.queryKey).toBe('15m:key42');
  await act(async () => root.render(<OverviewRealtimePanel realtime={response} loading window="15m" onWindowChange={vi.fn()} isDark isMobile={false} />));
  expect(speedProps.loading).toBe(false);
});

it('waits for the matching window response before giving candidates to a new speed query', async () => {
  await act(async () => root.render(<OverviewRealtimePanel realtime={response} loading={false} window="30m" onWindowChange={vi.fn()} isDark={false} isMobile={false} modelSpeedQueryKey="30m:key42" />));
  expect(speedProps.queryKey).toBe('30m:key42');
  expect(speedProps.speed).toBeUndefined();
  expect(speedProps.loading).toBe(true);
  const current = { ...response, window: '30m' as const };
  await act(async () => root.render(<OverviewRealtimePanel realtime={current} loading={false} window="30m" onWindowChange={vi.fn()} isDark={false} isMobile={false} modelSpeedQueryKey="30m:key42" />));
  expect(speedProps.speed).toBe(current.model_speed);
  expect(speedProps.loading).toBe(false);
});
