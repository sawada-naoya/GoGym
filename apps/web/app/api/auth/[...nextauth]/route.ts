import NextAuth, { NextAuthOptions, Session, User } from "next-auth";
import Credentials from "next-auth/providers/credentials";
import type { JWT } from "next-auth/jwt";
import type { AdapterUser } from "next-auth/adapters";

type LoginResponse = {
  user: { id: string | number; name: string; email: string };
  access_token: string;
  expires_in?: number;
};

const API_BASE = process.env.API_INTERNAL_URL ?? process.env.NEXT_PUBLIC_API_URL; // サーバ側用に内部URLがあれば優先

const authOptions: NextAuthOptions = {
  pages: {
    signIn: "/login",
    error: "/login",
  },

  session: { strategy: "jwt" },

  providers: [
    Credentials({
      name: "Email & Password",
      credentials: {
        email: { label: "Email", type: "email" },
        password: { label: "Password", type: "password" },
      },
      async authorize(credentials) {
        const email = credentials?.email?.trim();
        const password = credentials?.password;
        if (!email || !password) return null;

        const controller = new AbortController();
        const timeout = setTimeout(() => controller.abort(), 10_000);
        try {
          const r = await fetch(`${API_BASE}/api/v1/sessions/login`, {
            method: "POST",
            headers: { "content-type": "application/json" },
            body: JSON.stringify({ email, password }),
            cache: "no-store",
            signal: controller.signal,
          });
          clearTimeout(timeout);

          if (!r.ok) return null;

          const data = (await r.json()) as LoginResponse;
          const { user, access_token, expires_in } = data ?? {};
          if (!user?.id || !user?.name || !user?.email || !access_token) return null;

          return {
            id: String(user.id),
            name: user.name,
            email: user.email,
            accessToken: access_token,
            accessTokenExpiresAt: expires_in ? Date.now() + expires_in * 1000 : undefined,
          } as unknown as User;
        } catch {
          return null;
        } finally {
          clearTimeout(timeout);
        }
      },
    }),
  ],

  callbacks: {
    async jwt({ token, user }: { token: JWT; user?: User | AdapterUser }) {
      if (user) {
        // 初回ログイン時
        (token as any).accessToken = (user as any).accessToken;
        (token as any).accessTokenExpiresAt = (user as any).accessTokenExpiresAt;
        token.sub = user.id as string;
        token.name = user.name;
        token.email = user.email as string;
      }
      return token;
    },

    async session({ session, token }: { session: Session; token: JWT }) {
      (session as any).user = {
        id: token.sub,
        name: token.name,
        email: token.email,
      };
      (session as any).accessToken = (token as any).accessToken;
      (session as any).accessTokenExpiresAt = (token as any).accessTokenExpiresAt;
      return session;
    },

    async redirect({ url, baseUrl }) {
      try {
        if (url.startsWith("/")) return `${baseUrl}${url}`;
        const u = new URL(url);
        const base = new URL(baseUrl);

        if (u.origin === base.origin) return url;

        return `${baseUrl}/`;
      } catch {
        return baseUrl;
      }
    },
  },
};

const handler = NextAuth(authOptions);
export { handler as GET, handler as POST };
