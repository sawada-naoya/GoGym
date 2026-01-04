// server/accessToken.ts
import { headers } from "next/headers";
import { getToken } from "next-auth/jwt";

/**
 * サーバー側で accessToken を取得する
 * - Server Components / Server Actions / Route Handlers で使用
 * - sessionではなくJWTトークンから直接取得することで、クライアント露出を防ぐ
 */
export const getServerAccessToken = async (): Promise<string | null> => {
  const secret = process.env.NEXTAUTH_SECRET;
  if (!secret) return null;

  const headersList = await headers();
  const cookie = headersList.get("cookie") ?? "";

  const token = await getToken({
    req: { headers: { cookie } } as any,
    secret,
  });

  return (token?.accessToken as string) ?? null;
};
