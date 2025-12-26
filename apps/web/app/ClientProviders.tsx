"use client";

import { useEffect } from "react";
import { usePathname } from "next/navigation";
import { SessionProvider } from "next-auth/react";
import { I18nextProvider } from "react-i18next";
import { BannerProvider, useBanner } from "@/components/Banner";
import i18n from "@/lib/i18n/client";

const FlashMessageHandler = () => {
  const { success, error } = useBanner();
  const pathname = usePathname();

  useEffect(() => {
    const raw = sessionStorage.getItem("flash");
    if (!raw) return;

    try {
      const data = JSON.parse(raw) as {
        variant?: "success" | "error" | "info" | "warning";
        message?: string;
      };
      sessionStorage.removeItem("flash");

      if (data?.message) {
        if (data.variant === "error") {
          error(data.message);
        } else {
          success(data.message);
        }
      }
    } catch (e) {
      console.error("Failed to parse flash message:", e);
      sessionStorage.removeItem("flash");
    }
  }, [pathname, success, error]);

  return null;
};

const ClientProviders = ({ children }: { children: React.ReactNode }) => (
  <I18nextProvider i18n={i18n}>
    <SessionProvider>
      <BannerProvider>
        <FlashMessageHandler />
        {children}
      </BannerProvider>
    </SessionProvider>
  </I18nextProvider>
);

export default ClientProviders;
