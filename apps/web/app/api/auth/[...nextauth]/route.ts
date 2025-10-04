import NextAuth, { NextAuthOptions, Session, User } from "next-auth";
import Credentials from "next-auth/providers/credentials";
import type { JWT } from "next-auth/jwt";
import type { AdapterUser } from "next-auth/adapters";

type LoginResponse = {
  user: { id: string | number; name: string; email: string };
  access_token: string;
  expires_in?: number;
};

const API_BASE = process.env.NEXT_PUBLIC_API_URL;

const authOptions: NextAuthOptions = {
  providers: [
    Credentials({
      name: "Email & Password",
      credentials: { email: {}, password: {} },
      async authorize(credentials) {
        const r = await fetch(`${API_BASE}/api/v1/sessions/login`, {
          method: "POST",
          headers: { "content-type": "application/json" },
          body: JSON.stringify({
            email: credentials?.email,
            password: credentials?.password,
          }),
          cache: "no-store",
        });
        if (!r.ok) {
          return null;
        }

        const data = (await r.json()) as LoginResponse;
        const { user, access_token } = data ?? {};
        if (!user?.id || !user?.name || !user?.email || !access_token) return null;

        return {
          id: String(user.id),
          name: user.name,
          email: user.email,
          accessToken: access_token,
        } as any;
      },
    }),
  ],

  session: { strategy: "jwt" },

  callbacks: {
    async jwt({ token, user }: { token: JWT; user?: User | AdapterUser }) {
      if (user) {
        (token as any).accessToken = (user as any).accessToken;
      }
      return token;
    },

    async session({ session, token }: { session: Session; token: JWT }) {
      (session as any).accessToken = (token as any).accessToken;
      return session;
    },

    async redirect({ url, baseUrl }) {
      if (url.startsWith("/")) return `${baseUrl}${url}`;
      try {
        if (new URL(url).origin === baseUrl) return url;
      } catch {}
      return `${baseUrl}/login`;
    },
  },

  pages: {
    signIn: "/",
    error: "/login",
  },
};

const handler = NextAuth(authOptions);
export { handler as GET, handler as POST };
