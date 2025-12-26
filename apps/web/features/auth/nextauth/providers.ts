import type { User } from "next-auth";
import type { LoginResponse } from "../../../types/session";

const API_BASE =
  process.env.API_INTERNAL_URL ?? process.env.NEXT_PUBLIC_API_URL;

export const buildCredentialsProviderAuthorize = () => {
  return async (
    credentials: Partial<Record<"email" | "password", unknown>>,
    _request: Request,
  ) => {
    const email = credentials?.email?.toString().trim();
    const password = credentials?.password?.toString();
    if (!email || !password || !API_BASE) return null;

    try {
      const res = await fetch(`${API_BASE}/api/v1/sessions/login`, {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ email, password }),
        cache: "no-store",
      });
      if (!res.ok) return null;

      const data = (await res.json()) as LoginResponse;
      const { user, access_token, refresh_token, expires_in } = data ?? {};
      if (
        !user?.id ||
        !user?.name ||
        !user?.email ||
        !access_token ||
        !refresh_token
      ) {
        return null;
      }

      return {
        id: String(user.id),
        name: user.name,
        email: user.email,
        accessToken: access_token,
        refreshToken: refresh_token,
        expiresIn: expires_in,
      } as User;
    } catch {
      return null;
    }
  };
};
