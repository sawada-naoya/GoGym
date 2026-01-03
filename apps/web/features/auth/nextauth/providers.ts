import type { User } from "next-auth";
import type { LoginResponse } from "../../../types/session";

const API_BASE = process.env.NEXT_PUBLIC_API_URL;

export const buildCredentialsProviderAuthorize = () => {
  return async (credentials: Partial<Record<"email" | "password", unknown>>, _request: Request) => {
    const email = credentials?.email?.toString().trim();
    const password = credentials?.password?.toString();
    console.log("[DEBUG Login] Attempting login for:", email);
    console.log("[DEBUG Login] API_BASE:", API_BASE);

    if (!email || !password || !API_BASE) {
      console.error("[DEBUG Login] Missing credentials or API_BASE");
      return null;
    }

    try {
      const res = await fetch(`${API_BASE}/api/v1/sessions/login`, {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ email, password }),
        cache: "no-store",
      });

      console.log("[DEBUG Login] Response status:", res.status);
      if (!res.ok) {
        console.error("[DEBUG Login] Login failed with status:", res.status);
        return null;
      }

      const data = (await res.json()) as LoginResponse;
      const { user, access_token, refresh_token, expires_in } = data ?? {};

      console.log("[DEBUG Login] Response data:", {
        hasUser: !!user,
        hasAccessToken: !!access_token,
        hasRefreshToken: !!refresh_token,
        expiresIn: expires_in,
      });

      if (!user?.id || !user?.name || !user?.email || !access_token || !refresh_token) {
        console.error("[DEBUG Login] Missing required fields in response");
        return null;
      }

      console.log("[DEBUG Login] Login successful for user:", user.id);
      return {
        id: String(user.id),
        name: user.name,
        email: user.email,
        accessToken: access_token,
        refreshToken: refresh_token,
        expiresIn: expires_in,
      } as User;
    } catch (error) {
      console.error("[DEBUG Login] Login error:", error);
      return null;
    }
  };
};
