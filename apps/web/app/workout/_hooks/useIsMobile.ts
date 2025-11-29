import { useState, useEffect } from "react";

/**
 * 画面サイズがモバイル（768px未満）かどうかを判定するフック
 * @param breakpoint ブレークポイント（デフォルト: 768px）
 * @returns モバイルかどうかのboolean値
 */
export const useIsMobile = (breakpoint: number = 768): boolean => {
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < breakpoint);
    };

    checkMobile();
    window.addEventListener("resize", checkMobile);
    return () => window.removeEventListener("resize", checkMobile);
  }, [breakpoint]);

  return isMobile;
};
