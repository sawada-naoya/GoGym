export type BannerVariant = "success" | "error" | "info" | "warning";

export interface BannerProps {
  variant: BannerVariant;
  message: string | React.ReactNode;
  onClose?: () => void;
}
