// auth.ts
import NextAuth from "next-auth";
import Credentials from "next-auth/providers/credentials";
import { POST } from "@/lib/api";
import { LoginResponse } from "@/types/session";

export const { auth, signIn, signOut } = NextAuth({
  providers: [
    Credentials({
      name: "Email & Password",
      credentials: { email: {}, password: {} },
      async authorize(credentials) {
        // GoのログインAPIを叩く
        const res = await POST<LoginResponse>("/sessions/login", {
          body: {
            email: credentials?.email,
            password: credentials?.password,
          },
        });
        if (!res.ok || !res.data) return null;

        const { user, access_token } = res.data;
        if (!user?.id || !user?.name || !user?.email || !access_token) return null;

        // next-auth の“User”として返す（追加フィールドは any でOK）
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
    async jwt({ token, user }) {
      // サインイン直後だけ user が来る
      if (user) token.accessToken = (user as any).accessToken;
      return token;
    },
    async session({ session, token }) {
      (session as any).accessToken = token.accessToken;
      return session;
    },
  },

  pages: { signIn: "/login" },
});
