import { Banner } from "./Banner";
import type { BannerProps } from "./Banner.types";
export const ErrorBanner = (props: Omit<BannerProps, "variant">) => <Banner variant="error" {...props} />;
