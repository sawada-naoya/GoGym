import { Banner } from "./Banner";
import { BannerProps } from "./Banner.types";

export const SuccessBanner = (props: Omit<BannerProps, "variant">) => <Banner variant="success" {...props} />;
