// app/api/auth/[...nextauth]/route.ts
import NextAuth, { type NextAuthOptions, type Session, type User } from "next-auth";
import Credentials from "next-auth/providers/credentials";
import type { JWT } from "next-auth/jwt";
import type { AdapterUser } from "next-auth/adapters";

type LoginResponse = {
  user: { id: string | number; name: string; email: string };
  access_token: string;
  expires_in?: number;
};

const API_BASE = process.env.API_INTERNAL_URL ?? process.env.NEXT_PUBLIC_API_URL;

export const authOptions: NextAuthOptions = {
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
        const email = credentials?.email?.trim();
        const password = credentials?.password;
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
          } as unknown as User;
        } catch {
          return null;
        }
      },
    }),
  ],

  callbacks: {
    jwt: async ({ token, user }: { token: JWT; user?: User | AdapterUser }) => {
      if (user) {
        // accessToken は JWT トークン内にのみ保存（サーバー側のみ）
        (token as any).accessToken = (user as any).accessToken;
        token.sub = user.id as string;
        token.name = user.name;
        token.email = user.email as string;
      }
      return token;
    },

    session: async ({ session, token }: { session: Session; token: JWT }) => {
      // クライアント側に送信されるセッションにはaccessTokenを含めない
      (session as any).user = {
        id: token.sub,
        name: token.name,
        email: token.email,
      };
      // ⚠️ セキュリティ: accessToken はクライアント側に送信しない
      // サーバー側で必要な場合は getServerSession() から取得する
      return session;
    },

    // 最小：相対パスは許可、他はベースへ戻す
    redirect: async ({ url, baseUrl }) => (url.startsWith("/") ? `${baseUrl}${url}` : baseUrl),
  },
};

const handler = NextAuth(authOptions);
export { handler as GET, handler as POST };
