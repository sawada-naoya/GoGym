import NextAuth from "next-auth";
import Credentials from "next-auth/providers/credentials";
import type { User } from "next-auth";

type LoginResponse = {
  user: { id: string | number; name: string; email: string };
  access_token: string;
  expires_in?: number;
};

const API_BASE = process.env.API_INTERNAL_URL ?? process.env.NEXT_PUBLIC_API_URL;

export const { handlers, signIn, signOut, auth } = NextAuth({
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
          const { user, access_token } = data ?? {};
          if (!user?.id || !user?.name || !user?.email || !access_token) return null;

          return {
            id: String(user.id),
            name: user.name,
            email: user.email,
            accessToken: access_token,
          } as User;
        } catch {
          return null;
        }
      },
    }),
  ],

  callbacks: {
    jwt: async ({ token, user }) => {
      if (user) {
        // accessToken は JWT トークン内にのみ保存（サーバー側のみ）
        token.accessToken = (user as any).accessToken;
        token.sub = user.id as string;
      }
      return token;
    },

    session: async ({ session, token }) => {
      // サーバーサイドで auth() を使用する場合は accessToken を含める
      // クライアントサイドには送信されない（Server Components/Server Actions のみ）
      if (session.user) {
        session.user.id = token.sub as string;
        session.user.accessToken = token.accessToken as string;
      }
      return session;
    },
  },
});
