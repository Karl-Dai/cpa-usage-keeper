import { useCallback, useEffect, useRef, type RefObject } from 'react';
import { createPortal } from 'react-dom';
import { useTranslation } from 'react-i18next';
import { IconX } from '@/components/ui/icons';
import { useAnchorPosition } from '@/hooks/useAnchorPosition';
import panelStyles from '../analysis/AnalysisPanel.module.scss';
import styles from './ModelSpeedChart.module.scss';

export function ModelSpeedTooltip({ canvasRef, x, y, touch, title, rows, dismiss, cancelClose, scheduleClose }: {
  canvasRef: RefObject<HTMLCanvasElement | null>;
  x: number;
  y: number;
  touch: boolean;
  title: string;
  rows: { model: string; color: string; detail: string }[];
  dismiss: () => void;
  cancelClose: () => void;
  scheduleClose: () => void;
}) {
  const { t } = useTranslation();
  const panelRef = useRef<HTMLDivElement>(null);
  const update = useCallback((rect: DOMRect) => {
    const panel = panelRef.current;
    if (!panel) return;
    const viewport = window.visualViewport;
    const width = viewport?.width ?? window.innerWidth;
    const height = viewport?.height ?? window.innerHeight;
    const leftEdge = (viewport?.offsetLeft ?? 0) + 8;
    const topEdge = (viewport?.offsetTop ?? 0) + 8;
    panel.style.width = `${Math.min(360, width - 16)}px`;
    panel.style.maxHeight = `${Math.min(380, touch ? height * 0.6 : height - 16)}px`;
    const bounds = panel.getBoundingClientRect();
    const rightEdge = leftEdge + width - 16;
    const bottomEdge = topEdge + height - 16;
    const cursorX = rect.left + x;
    const cursorY = rect.top + y;
    const preferredLeft = cursorX + 10 + bounds.width <= rightEdge ? cursorX + 10 : cursorX - bounds.width - 10;
    panel.style.left = `${window.scrollX + (touch ? leftEdge : Math.max(leftEdge, Math.min(preferredLeft, rightEdge - bounds.width)))}px`;
    panel.style.top = `${window.scrollY + (touch ? bottomEdge - bounds.height : Math.max(topEdge, Math.min(cursorY - bounds.height / 2, bottomEdge - bounds.height)))}px`;
  }, [touch, x, y]);
  useAnchorPosition(true, canvasRef, update);
  useEffect(() => {
    const outside = (event: PointerEvent) => {
      if (event.target instanceof Node && event.target !== canvasRef.current && !panelRef.current?.contains(event.target)) dismiss();
    };
    const escape = (event: KeyboardEvent) => { if (event.key === 'Escape') dismiss(); };
    document.addEventListener('pointerdown', outside);
    document.addEventListener('keydown', escape);
    return () => { document.removeEventListener('pointerdown', outside); document.removeEventListener('keydown', escape); };
  }, [canvasRef, dismiss]);

  return createPortal(<div ref={panelRef} role="tooltip" className={`${panelStyles.topModelsTooltip} ${styles.tooltip}`} data-model-speed-tooltip
    onPointerEnter={cancelClose} onPointerLeave={scheduleClose}>
    <div className={panelStyles.topModelsTooltipHeader}>
      <strong>{title}</strong>
      <button type="button" aria-label={t('common.close')} onClick={dismiss}><IconX size={16} /></button>
    </div>
    <ol className={panelStyles.topModelsTooltipRows}>
      {rows.map((row) => <li key={row.model}>
        <span className={panelStyles.topModelsTooltipDot} style={{ backgroundColor: row.color }} aria-hidden="true" />
        <span><strong>{row.model}</strong><br />{row.detail}</span>
      </li>)}
    </ol>
  </div>, document.body);
}
