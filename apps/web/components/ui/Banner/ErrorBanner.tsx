import { Banner } from "./Banner";
import { ApiError, getFormError } from "../../../lib/tokenStore";

interface ErrorBannerProps {
  message?: string | React.ReactNode;
  error?: ApiError | null;
  onClose?: () => void;
}

export const ErrorBanner = ({ message, error, onClose }: ErrorBannerProps) => {
  const displayMessage = message || getFormError(error);
  if (!displayMessage) return null;

  return <Banner variant="error" message={displayMessage} onClose={onClose} />;
};
