import NextAuth from "next-auth";
import Credentials from "next-auth/providers/credentials";
import type { User } from "next-auth";
import type { LoginResponse } from "../../types/session";
import type { RefreshResponse } from "../../types/refresh";

const API_BASE = process.env.NEXT_PUBLIC_API_URL;
const secret = process.env.AUTH_SECRET;

if (!secret && process.env.NODE_ENV === "production") {
  throw new Error("AUTH_SECRET is not set in production environment");
}

/**
 * リフレッシュトークンを使って新しいアクセストークンを取得
 */
export const refreshAccessToken = async (refreshToken: string): Promise<RefreshResponse | null> => {
  try {
    const res = await fetch(`${API_BASE}/api/v1/sessions/refresh`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
      cache: "no-store",
    });

    if (!res.ok) {
      console.error("Failed to refresh token:", res.status);
      return null;
    }

    const data = (await res.json()) as RefreshResponse;
    return data;
  } catch (error) {
    console.error("Error refreshing token:", error);
    return null;
  }
};

/**
 * Credentials Provider の authorize 関数
 */
const buildCredentialsProviderAuthorize = () => {
  return async (credentials: Partial<Record<"email" | "password", unknown>>, _request: Request) => {
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
      if (!user?.id || !user?.name || !user?.email || !access_token || !refresh_token) {
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

/**
 * NextAuth設定
 */
export const { handlers, signIn, signOut, auth } = NextAuth({
  secret: secret,
  trustHost: true,

  pages: {
    signIn: "/",
    error: "/",
  },
  session: { strategy: "jwt" },

  providers: [
    Credentials({
      name: "Email & Password",
      credentials: {
        email: { label: "Email", type: "email" },
        password: { label: "Password", type: "password" },
      },
      authorize: buildCredentialsProviderAuthorize(),
    }),
  ],

  callbacks: {
    jwt: async ({ token, user }) => {
      if (user) {
        const expiresAt = Math.floor(Date.now() / 1000) + ((user as any).expiresIn || 900);

        return {
          ...token,
          accessToken: (user as any).accessToken,
          refreshToken: (user as any).refreshToken,
          expiresAt,
          sub: user.id as string,
        };
      }

      const now = Math.floor(Date.now() / 1000);
      const expiresAt = (token.expiresAt as number) || 0;

      if (now < expiresAt - 60) return token;

      const refreshed = await refreshAccessToken(token.refreshToken as string);
      if (!refreshed) return { ...token, error: "RefreshTokenError" };

      const newExpiresAt = Math.floor(Date.now() / 1000) + refreshed.expires_in;
      return {
        ...token,
        accessToken: refreshed.access_token,
        refreshToken: refreshed.refresh_token,
        expiresAt: newExpiresAt,
      };
    },

    session: async ({ session, token }) => {
      (session as any).authError = token.error ?? null;
      if (session.user) session.user.id = token.sub as string;
      return session;
    },
  },
});
