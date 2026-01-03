import NextAuth from "next-auth";
import Credentials from "next-auth/providers/credentials";

import { refreshAccessToken } from "./refresh";
import { buildCredentialsProviderAuthorize } from "./providers";

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
      authorize: buildCredentialsProviderAuthorize(),
    }),
  ],

  callbacks: {
    jwt: async ({ token, user }) => {
      if (user) {
        console.log("[DEBUG JWT] New user login");
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
      console.log("[DEBUG JWT] Token check - now:", now, "expiresAt:", expiresAt, "timeLeft:", expiresAt - now);

      if (now < expiresAt - 60) {
        console.log("[DEBUG JWT] Token still valid, skipping refresh");
        return token;
      }

      console.log("[DEBUG JWT] Token expired, attempting refresh");
      const refreshed = await refreshAccessToken(token.refreshToken as string);
      if (!refreshed) {
        console.error("[DEBUG JWT] Refresh failed, setting RefreshTokenError");
        return { ...token, error: "RefreshTokenError" };
      }

      console.log("[DEBUG JWT] Refresh successful");
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
      console.log("[DEBUG Session] authError:", token.error);
      return session;
    },
  },
});
