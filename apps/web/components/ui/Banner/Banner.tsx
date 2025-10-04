import { useEffect } from "react";
import { BannerProps, BannerVariant } from "./Banner.types";

const colors: Record<BannerVariant, { bg: string; border: string; text: string; icon: string; ring: string }> = {
  success: { bg: "bg-green-50", border: "border-green-300", text: "text-green-800", icon: "text-green-400", ring: "focus:ring-green-600" },
  error: { bg: "bg-red-50", border: "border-red-300", text: "text-red-800", icon: "text-red-400", ring: "focus:ring-red-600" },
  info: { bg: "bg-blue-50", border: "border-blue-300", text: "text-blue-800", icon: "text-blue-400", ring: "focus:ring-blue-600" },
  warning: { bg: "bg-yellow-50", border: "border-yellow-300", text: "text-yellow-800", icon: "text-yellow-400", ring: "focus:ring-yellow-600" },
};

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

export const Banner = ({ variant, message, title, dismissible, autoHideMs, onClose, icon = true, className = "", "data-testid": testId }: BannerProps) => {
  const { bg, border, text, icon: iconColor, ring } = colors[variant];

  const isDismissible = dismissible ?? !!onClose;

  useEffect(() => {
    if (!autoHideMs || !onClose) return;
    const t = setTimeout(onClose, autoHideMs);
    return () => clearTimeout(t);
  }, [autoHideMs, onClose]);

  // アクセシビリティ：エラー/警告は assertive、それ以外は polite
  const ariaLive = variant === "error" || variant === "warning" ? "assertive" : "polite";
  const role = variant === "error" || variant === "warning" ? "alert" : "status";

  return (
    <div className={`mb-4 rounded-md border ${border} ${bg} p-4 ${className}`} role={role} aria-live={ariaLive} data-testid={testId}>
      <div className="flex">
        {icon !== false && (
          <div className="flex-shrink-0">
            <div className={iconColor}>{icon === true ? defaultIcons[variant] : icon}</div>
          </div>
        )}

        <div className={`ml-3 flex-1 ${icon === false ? "" : "mt-0.5"}`}>
          {title && <div className={`text-sm font-semibold ${text}`}>{title}</div>}
          <div className={`text-sm ${text}`}>{message}</div>
        </div>

        {isDismissible && onClose && (
          <div className="ml-auto pl-3">
            <button type="button" onClick={onClose} className={`inline-flex rounded-md p-1.5 ${colors[variant].icon} hover:opacity-80 focus:outline-none focus:ring-2 ${ring}`} aria-label="閉じる">
              <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
              </svg>
            </button>
          </div>
        )}
      </div>
    </div>
  );
};
