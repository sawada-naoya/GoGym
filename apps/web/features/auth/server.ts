// server/accessToken.ts
import { cookies } from "next/headers";
import { getToken } from "next-auth/jwt";

export const getServerAccessToken = async () => {
  const secret = process.env.AUTH_SECRET;
  if (!secret) return null;

  const cookieStore = await cookies();
  const cookieHeader = cookieStore
    .getAll()
    .map((c) => `${c.name}=${c.value}`)
    .join("; ");

  const token = await getToken({ req: { headers: { cookie: cookieHeader } } as any, secret });
  return (token?.accessToken as string) ?? null;
};
