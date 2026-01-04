// server/accessToken.ts
import { headers } from "next/headers";
import { getToken } from "next-auth/jwt";

type AccessTokenDiag = {
  accessToken: string | null;
  reason: "OK" | "NO_SECRET" | "NO_COOKIE_HEADER" | "GETTOKEN_NULL" | "GETTOKEN_THROW";
  meta: {
    hasSecret: boolean;
    cookieHeaderLength: number;
    hasNextAuthCookieName: boolean;
    // ここは値を出さない。名前と存在だけ。
    cookieNamesSample: string[];
  };
  error?: { name?: string; message?: string };
};

const parseCookieNames = (cookieHeader: string) => {
  // cookie 値は出さない（名前だけ）
  return cookieHeader
    .split(";")
    .map((p) => p.trim().split("=")[0])
    .filter(Boolean);
};

export const getServerAccessTokenWithDiag = async (): Promise<AccessTokenDiag> => {
  const secret = process.env.NEXTAUTH_SECRET;
  if (!secret) {
    return {
      accessToken: null,
      reason: "NO_SECRET",
      meta: {
        hasSecret: false,
        cookieHeaderLength: 0,
        hasNextAuthCookieName: false,
        cookieNamesSample: [],
      },
    };
  }

  const headersList = await headers();
  const cookieHeader = headersList.get("cookie") ?? "";
  const cookieNames = parseCookieNames(cookieHeader);

  const hasNextAuthCookieName = cookieNames.includes("next-auth.session-token") || cookieNames.includes("__Secure-next-auth.session-token") || cookieNames.includes("__Host-next-auth.session-token");

  if (!cookieHeader) {
    return {
      accessToken: null,
      reason: "NO_COOKIE_HEADER",
      meta: {
        hasSecret: true,
        cookieHeaderLength: 0,
        hasNextAuthCookieName,
        cookieNamesSample: cookieNames.slice(0, 20),
      },
    };
  }

  try {
    const token = await getToken({
      req: { headers: { cookie: cookieHeader } } as any,
      secret,
    });

    if (!token) {
      return {
        accessToken: null,
        reason: "GETTOKEN_NULL",
        meta: {
          hasSecret: true,
          cookieHeaderLength: cookieHeader.length,
          hasNextAuthCookieName,
          cookieNamesSample: cookieNames.slice(0, 20),
        },
      };
    }

    return {
      accessToken: (token?.accessToken as string) ?? null,
      reason: "OK",
      meta: {
        hasSecret: true,
        cookieHeaderLength: cookieHeader.length,
        hasNextAuthCookieName,
        cookieNamesSample: cookieNames.slice(0, 20),
      },
    };
  } catch (e: any) {
    return {
      accessToken: null,
      reason: "GETTOKEN_THROW",
      meta: {
        hasSecret: true,
        cookieHeaderLength: cookieHeader.length,
        hasNextAuthCookieName,
        cookieNamesSample: cookieNames.slice(0, 20),
      },
      error: { name: e?.name, message: e?.message },
    };
  }
};

// 本番利用用（従来互換）
export const getServerAccessToken = async (): Promise<string | null> => {
  const r = await getServerAccessTokenWithDiag();
  return r.accessToken;
};
