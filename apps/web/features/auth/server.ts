import { cookies } from "next/headers";
import { getToken } from "next-auth/jwt";

/**
 * 環境に応じたセッションクッキー名を取得
 * - Production (HTTPS): __Secure-authjs.session-token
 * - Development (HTTP): authjs.session-token
 */
const getSessionCookieName = (): string => {
  const isProduction = process.env.NODE_ENV === "production";
  const useSecureCookies = process.env.AUTH_URL?.startsWith("https://") ?? isProduction;

  if (useSecureCookies) {
    return "__Secure-authjs.session-token";
  }
  return "authjs.session-token";
};

export const getServerAccessToken = async (): Promise<string | null> => {
  const secret = process.env.AUTH_SECRET;
  if (!secret) return null;

  const store = await cookies();
  const allCookies = store.getAll();
  const cookieHeader = allCookies.map((c) => `${c.name}=${c.value}`).join("; ");

  const expectedCookieName = getSessionCookieName();
  const sessionCookie = store.get(expectedCookieName);

  if (!sessionCookie) {
    // フォールバック: すべてのセッションクッキーを試す
    const fallbackCookie = store.get("__Secure-authjs.session-token") ?? store.get("__Host-authjs.session-token") ?? store.get("authjs.session-token");

    if (!fallbackCookie) {
      return null;
    }
  }

  const token = await getToken({
    req: { headers: { cookie: cookieHeader } } as any,
    secret,
    cookieName: expectedCookieName,
  });

  if (!token) return null;
  if (!token.accessToken) return null;

  return token.accessToken as string;
};
