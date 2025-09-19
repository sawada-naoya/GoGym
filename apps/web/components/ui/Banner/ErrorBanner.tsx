import { Banner } from "./Banner";
import { BannerProps } from "./Banner.types";
import { ApiError, getFormError } from "../../../lib/errors";

interface ErrorBannerProps {
  message?: string | React.ReactNode;
  error?: ApiError | null;
  onClose?: () => void;
}

export function ErrorBanner({ message, error, onClose }: ErrorBannerProps) {
  const displayMessage = message || getFormError(error);
  if (!displayMessage) return null;

  return <Banner variant="error" message={displayMessage} onClose={onClose} />;
}
