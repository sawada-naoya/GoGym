import { getToken } from "next-auth/jwt";
import { headers } from "next/headers";

/**
 * サーバー側でJWT全体を取得するヘルパー関数
 *
 * 使用可能な場所:
 * - Server Components
 * - Server Actions
 * - API Route Handlers
 *
 * @returns JWT token object or null
 */
export const getServerToken = async () => {
  const secret = process.env.NEXTAUTH_SECRET;
  if (!secret) return null;

  const headersList = await headers();
  const cookie = headersList.get("cookie") ?? "";
  const token = await getToken({
    req: { headers: { cookie } } as any,
    secret,
  });

  return token ?? null;
};

/**
 * サーバー側でJWTからaccessTokenを取得するヘルパー関数
 *
 * 使用可能な場所:
 * - Server Components
 * - Server Actions
 * - API Route Handlers
 *
 * @returns accessToken or null if not found
 */
export const getServerAccessToken = async (): Promise<string | null> => {
  const token = await getServerToken();
  return (token?.accessToken as string) || null;
};
