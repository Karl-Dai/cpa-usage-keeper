import { useCallback, useRef } from 'react';
import type { Chart } from 'chart.js';
import { useMediaQuery } from '@/hooks/useMediaQuery';

export const MODEL_SPEED_ANIMATION_MS = 500;
export type ModelSpeedRefreshSignal = { source: object; revision: number };

export function useModelSpeedAnimation(queryKey: string, source: object | undefined, refreshSignal?: ModelSpeedRefreshSignal) {
  const reducedMotion = useMediaQuery('(prefers-reduced-motion: reduce)');
  const previous = useRef<{ chart: Chart<'scatter'>; queryKey: string; revision: number } | null>(null);
  const beforeDatasetsUpdate = useCallback((chart: Chart<'scatter'>) => {
    const last = previous.current;
    const revision = refreshSignal && refreshSignal.source === source ? refreshSignal.revision : null;
    const refreshed = revision !== null && revision !== last?.revision;
    const replay = !last || last.chart !== chart || last.queryKey !== queryKey || refreshed;
    previous.current = { chart, queryKey, revision: revision ?? last?.revision ?? 0 };
    // 布局完成后由原生 reset/update 重播数值轴；背景数据、选择和 resize 继续从现有坐标过渡。
    if (replay && !reducedMotion) { chart.stop(); chart.reset(); }
  }, [queryKey, source, refreshSignal, reducedMotion]);
  return { reducedMotion, beforeDatasetsUpdate };
}
