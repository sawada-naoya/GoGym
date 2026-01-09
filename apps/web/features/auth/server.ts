import { cookies } from "next/headers";
import { getToken } from "next-auth/jwt";

export const getServerAccessToken = async (): Promise<string | null> => {
  const secret = process.env.AUTH_SECRET;
  if (!secret) return null;

  const store = await cookies();
  const cookieHeader = store
    .getAll()
    .map((c) => `${c.name}=${c.value}`)
    .join("; ");

  const token = await getToken({
    req: { headers: { cookie: cookieHeader } } as any,
    secret,
    cookieName: "__Secure-authjs.session-token",
  });

  return (token?.accessToken as string) ?? null;
};
