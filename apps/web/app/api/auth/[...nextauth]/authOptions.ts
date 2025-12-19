import NextAuth from "next-auth";
import Credentials from "next-auth/providers/credentials";
import type { User } from "next-auth";

type LoginResponse = {
  user: { id: string | number; name: string; email: string };
  access_token: string;
  refresh_token: string;
  expires_in: number;
};

type RefreshResponse = {
  user: { id: string; name: string; email: string };
  access_token: string;
  refresh_token: string;
  expires_in: number;
};

const API_BASE =
  process.env.API_INTERNAL_URL ?? process.env.NEXT_PUBLIC_API_URL;

/**
 * リフレッシュトークンを使って新しいアクセストークンを取得
 */
const refreshAccessToken = async (
  refreshToken: string,
): Promise<RefreshResponse | null> => {
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

export const { handlers, signIn, signOut, auth } = NextAuth({
  pages: {
    signIn: "/",
    error: "/",
  },
  session: { strategy: "jwt" },

  // NextAuth に「メール/パスワードでログインする」入口を提供
  providers: [
    Credentials({
      name: "Email & Password",
      credentials: {
        email: { label: "Email", type: "email" },
        password: { label: "Password", type: "password" },
      },
      // authorize が Go API (/api/v1/sessions/login)に問い合わせて認証を行う
      authorize: async (credentials) => {
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
          )
            return null;

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
      },
    }),
  ],

  // NextAuth のセッションを DB ではなく JWT で管理するための設定
  callbacks: {
    jwt: async ({ token, user }) => {
      // 初回ログイン時: userオブジェクトからトークン情報を保存
      if (user) {
        const expiresAt =
          Math.floor(Date.now() / 1000) + ((user as any).expiresIn || 900); // 15分 = 900秒
        return {
          ...token,
          accessToken: (user as any).accessToken,
          refreshToken: (user as any).refreshToken,
          expiresAt,
          sub: user.id as string,
        };
      }

      // トークンの有効期限をチェック（期限切れ1分前にリフレッシュ）
      const now = Math.floor(Date.now() / 1000);
      const expiresAt = (token.expiresAt as number) || 0;

      // まだ有効期限内（1分以上余裕がある）場合はそのまま返す
      if (now < expiresAt - 60) {
        return token;
      }

      // トークンをリフレッシュ
      console.log("Refreshing access token...");
      const refreshed = await refreshAccessToken(token.refreshToken as string);

      if (!refreshed) {
        // リフレッシュ失敗 → ログアウトさせるためにエラートークンを返す
        console.error("Failed to refresh token, marking for logout");
        return { ...token, error: "RefreshTokenError" };
      }

      // 新しいトークン情報で更新
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

      if (session.user) {
        session.user.id = token.sub as string;
      }

      return session;
    },
  },
});
