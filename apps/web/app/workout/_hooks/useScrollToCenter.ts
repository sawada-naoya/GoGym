import { useEffect, RefObject } from "react";

/**
 * 選択された要素をコンテナの中央にスクロールするフック
 *
 * @param containerRef スクロールコンテナのref
 * @param selectedElementRef 選択された要素のref
 */
export const useScrollToCenter = <T extends HTMLElement = HTMLElement, U extends HTMLElement = HTMLElement>(
  containerRef: RefObject<T | null>,
  selectedElementRef: RefObject<U | null>
) => {
  useEffect(() => {
    if (containerRef.current && selectedElementRef.current) {
      const container = containerRef.current;
      const element = selectedElementRef.current;

      const containerWidth = container.clientWidth;
      const elementLeft = element.offsetLeft;
      const elementWidth = element.clientWidth;

      // 要素の中心位置 - コンテナの中心位置 = スクロール位置
      const scrollPosition = elementLeft - containerWidth / 2 + elementWidth / 2;

      container.scrollLeft = scrollPosition;
    }
  }, []); // 初回のみ実行
};
