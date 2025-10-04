export type BannerVariant = "success" | "error" | "info" | "warning";

export interface BannerProps {
  variant: BannerVariant;
  message: string | React.ReactNode;
  title?: string;
  dismissible?: boolean; // × を出すか（onClose がある場合は自動的に true 扱いでもOK）
  autoHideMs?: number; // 自動で閉じる（onClose 指定時のみ動作）
  onClose?: () => void;
  icon?: boolean | React.ReactNode; // false で非表示、ReactNode を渡せば差し替え
  className?: string;
  "data-testid"?: string;
}
