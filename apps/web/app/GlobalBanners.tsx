// app/GlobalBanners.tsx
"use client";
import { useEffect } from "react";
import { usePathname } from "next/navigation";
import { useBanner } from "@/components/Banner";

const GlobalBanners = () => {
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
        // データが取得できたら表示
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

export default GlobalBanners;
