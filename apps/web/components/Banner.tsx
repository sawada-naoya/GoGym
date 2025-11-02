// components/ui/Banner.tsx
"use client";
import React, { createContext, useCallback, useContext, useMemo, useState } from "react";

/* types */
export type BannerVariant = "success" | "error" | "info" | "warning";
type BannerItem = { id: string; variant: BannerVariant; message: React.ReactNode; autoHideMs: number; icon?: boolean | React.ReactNode; className?: string };
type PushOptions = { autoHideMs?: number; icon?: BannerItem["icon"]; className?: string };
type BannerContextType = {
  items: BannerItem[];
  push: (variant: BannerVariant, message: React.ReactNode, opts?: PushOptions) => void;
};

const BannerCtx = createContext<BannerContextType | null>(null);

/* styles & icons（前と同じ） */
const colors = {
  success: { bg: "bg-green-50", border: "border-green-300", text: "text-green-800", icon: "text-green-400" },
  error: { bg: "bg-red-50", border: "border-red-300", text: "text-red-800", icon: "text-red-400" },
  info: { bg: "bg-blue-50", border: "border-blue-300", text: "text-blue-800", icon: "text-blue-400" },
  warning: { bg: "bg-yellow-50", border: "border-yellow-300", text: "text-yellow-800", icon: "text-yellow-400" },
} as const;

const defaultIcons: Record<BannerVariant, JSX.Element> = {
  success: (
    <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
    </svg>
  ),
  error: (
    <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
    </svg>
  ),
  info: (
    <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
      <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a.75.75 0 000 1.5h.253a.25.25 0 01.244.304l-.459 2.066A1.75 1.75 0 0010.747 15H11a.75.75 0 000-1.5h-.253a.25.25 0 01-.244-.304l.459-2.066A1.75 1.75 0 009.253 9H9z" clipRule="evenodd" />
    </svg>
  ),
  warning: (
    <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
      <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
    </svg>
  ),
};

/* Provider（全ツリーを包む） */
export const BannerProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [items, setItems] = useState<BannerItem[]>([]);

  const push = useCallback<BannerContextType["push"]>((variant, message, opts) => {
    const id = crypto.randomUUID();
    const autoHideMs = opts?.autoHideMs ?? 5000;
    const item: BannerItem = { id, variant, message, autoHideMs, icon: opts?.icon, className: opts?.className };
    setItems((prev) => [...prev, item]);
    window.setTimeout(() => setItems((prev) => prev.filter((x) => x.id !== id)), autoHideMs);
  }, []);

  const value = useMemo(() => ({ items, push }), [items, push]);
  return <BannerCtx.Provider value={value}>{children}</BannerCtx.Provider>;
};

/* Host（画面に描画する"置き場所"） */
export const BannerHost: React.FC<{ className?: string }> = ({ className }) => {
  const ctx = useContext(BannerCtx);
  if (!ctx) throw new Error("<BannerHost /> must be used under <BannerProvider />.");

  return (
    <div className={`fixed top-20 left-1/2 -translate-x-1/2 z-40 w-full max-w-2xl px-4 space-y-2 ${className ?? ""}`}>
      {ctx.items.map(({ id, variant, message, autoHideMs, icon, className: itemClass }) => {
        const { bg, border, text, icon: iconColor } = colors[variant];
        const ariaLive = variant === "error" || variant === "warning" ? "assertive" : "polite";
        const role = variant === "error" || variant === "warning" ? "alert" : "status";

        // アイテム自身も auto-hide で消えるけど、Provider 側でリストからも消しているので
        // ここでは表示のみ（WYSIWYG）に徹する
        return (
          <div key={id} className={`mb-4 rounded-md border ${border} ${bg} p-4 ${itemClass ?? ""}`} role={role} aria-live={ariaLive}>
            <div className="flex">
              {icon !== false && (
                <div className="flex-shrink-0">
                  <div className={iconColor}>{icon === true ? defaultIcons[variant] : icon}</div>
                </div>
              )}
              <div className={`ml-3 flex-1 ${icon === false ? "" : "mt-0.5"}`}>
                <div className={`text-sm ${text}`}>{typeof message === "string" ? <span>{message}</span> : message}</div>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};

/* Hook（RHFっぽく使う） */
export const useBanner = () => {
  const ctx = useContext(BannerCtx);
  if (!ctx) throw new Error("useBanner must be used under <BannerProvider />.");
  const success = (message: React.ReactNode, opts?: PushOptions) => ctx.push("success", message, opts);
  const error = (message: React.ReactNode, opts?: PushOptions) => ctx.push("error", message, opts);
  const info = (message: React.ReactNode, opts?: PushOptions) => ctx.push("info", message, opts);
  const warning = (message: React.ReactNode, opts?: PushOptions) => ctx.push("warning", message, opts);
  return { success, error, info, warning };
};
