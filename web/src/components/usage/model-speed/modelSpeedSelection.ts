import { useCallback, useEffect, useMemo, useState } from 'react';

export const MODEL_SPEED_SELECTION_LIMIT = 5;
export type ModelSpeedCandidate = { model: string; total_tokens: number };
export type ModelSpeedView = 'analysis' | 'realtime';

export const modelSpeedStorageKey = (view: ModelSpeedView, scope: string) =>
  `cpa-usage-keeper-model-speed-v1:${view}:${scope}`;

export function sortModelSpeedCandidates(candidates: readonly ModelSpeedCandidate[]) {
  return [...candidates].sort((left, right) => right.total_tokens - left.total_tokens
    || (left.model < right.model ? -1 : left.model > right.model ? 1 : 0));
}

function loadSelection(key: string): { models: string[]; colors: Record<string, number> } | null {
  try {
    const raw = localStorage.getItem(key);
    if (!raw) return null;
    const value: unknown = JSON.parse(raw);
    if (!value || typeof value !== 'object' || !('version' in value) || value.version !== 1
      || !('models' in value) || !Array.isArray(value.models)
      || !value.models.every((model) => typeof model === 'string' && model.trim())) return null;
    const models = [...new Set(value.models as string[])].slice(0, MODEL_SPEED_SELECTION_LIMIT);
    const saved: Record<string, number> = {};
    if ('colors' in value && value.colors && typeof value.colors === 'object') {
      for (const [model, index] of Object.entries(value.colors)) {
        if (typeof index === 'number' && Number.isInteger(index) && index >= 0 && index < MODEL_SPEED_SELECTION_LIMIT) saved[model] = index;
      }
    }
    return { models, colors: assignColors(models, saved) };
  } catch { return null; }
}

function persistSelection(key: string, models: readonly string[], colors: Readonly<Record<string, number>>) {
  try { localStorage.setItem(key, JSON.stringify({ version: 1, models, colors })); }
  catch { /* 浏览器存储受限时继续保留本次选择。 */ }
}

function assignColors(models: readonly string[], previous: Readonly<Record<string, number>> = {}) {
  const colors: Record<string, number> = {};
  const used = new Set<number>();
  for (const model of models) {
    if (Object.hasOwn(previous, model) && !used.has(previous[model])) { colors[model] = previous[model]; used.add(previous[model]); }
  }
  for (const model of models) {
    if (Object.hasOwn(colors, model)) continue;
    const available = Array.from({ length: MODEL_SPEED_SELECTION_LIMIT }, (_, index) => index).find((index) => !used.has(index)) ?? 0;
    colors[model] = available;
    used.add(available);
  }
  return colors;
}

export function useModelSpeedSelection(view: ModelSpeedView, scope: string, candidates: readonly ModelSpeedCandidate[]) {
  const storageKey = modelSpeedStorageKey(view, scope);
  const sortedCandidates = useMemo(() => sortModelSpeedCandidates(candidates), [candidates]);
  const defaults = useMemo(() => sortedCandidates.slice(0, MODEL_SPEED_SELECTION_LIMIT).map((item) => item.model), [sortedCandidates]);
  const [state, setState] = useState(() => {
    const saved = loadSelection(storageKey);
    return { key: storageKey, models: saved?.models ?? null, colors: saved?.colors ?? {} };
  });
  // 切换 Viewer 时同步重置作用域；只在首次有候选数据时自动选默认，显式空选不会重新初始化。
  if (state.key !== storageKey) {
    const saved = loadSelection(storageKey);
    setState({ key: storageKey, models: saved?.models ?? null, colors: saved?.colors ?? {} });
  } else if (state.models === null && defaults.length > 0) {
    setState({ key: storageKey, models: defaults, colors: assignColors(defaults) });
  }
  useEffect(() => {
    if (state.key === storageKey && state.models !== null) persistSelection(storageKey, state.models, state.colors);
  }, [state, storageKey]);

  const choose = useCallback((models: string[]) => {
    const next = [...new Set(models)].slice(0, MODEL_SPEED_SELECTION_LIMIT);
    const colors = assignColors(next, state.key === storageKey ? state.colors : {});
    persistSelection(storageKey, next, colors);
    setState({ key: storageKey, models: next, colors });
  }, [state.colors, state.key, storageKey]);

  return { selected: state.models ?? [], colors: state.colors, candidates: sortedCandidates, choose, reset: () => choose(defaults) };
}
