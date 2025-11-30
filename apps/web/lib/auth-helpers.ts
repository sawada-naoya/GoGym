import { getToken } from "next-auth/jwt";
import { cookies } from "next/headers";

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
  try {
    const cookieStore = cookies();
    const token = await getToken({
      req: { cookies: cookieStore } as any,
      secret: process.env.NEXTAUTH_SECRET,
    });

    return (token?.accessToken as string) || null;
  } catch (error) {
    console.error("Failed to get access token:", error);
    return null;
  }
};

/**
 * サーバー側でJWT全体を取得するヘルパー関数
 *
 * @returns JWT token object or null
 */
export const getServerToken = async () => {
  try {
    const cookieStore = cookies();
    const token = await getToken({
      req: { cookies: cookieStore } as any,
      secret: process.env.NEXTAUTH_SECRET,
    });

    return token;
  } catch (error) {
    console.error("Failed to get token:", error);
    return null;
  }
};
