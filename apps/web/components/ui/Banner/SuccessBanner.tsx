import { Banner } from "./Banner";
import type { BannerProps } from "./Banner.types";
export const SuccessBanner = (props: Omit<BannerProps, "variant">) => <Banner variant="success" {...props} />;
