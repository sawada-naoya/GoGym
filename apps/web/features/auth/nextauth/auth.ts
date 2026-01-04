import NextAuth from "next-auth";
import Credentials from "next-auth/providers/credentials";

import { refreshAccessToken } from "./refresh";
import { buildCredentialsProviderAuthorize } from "./providers";

const secret = process.env.NEXTAUTH_SECRET;

if (!secret && process.env.NODE_ENV === "production") {
  throw new Error("NEXTAUTH_SECRET is not set in production environment");
}

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
      (session as any).expiresAt = token.expiresAt ?? null;
      return session;
    },
  },
});
