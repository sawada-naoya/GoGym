import { NextResponse } from "next/server";

export const dynamic = "force-dynamic";

export const GET = async () => {
  const authSecret = process.env.AUTH_SECRET ?? "";

  return NextResponse.json({
    hasSecret: Boolean(authSecret),
    length: authSecret.length,
    firstCharCode: authSecret.charCodeAt(0),
    lastCharCode: authSecret.charCodeAt(authSecret.length - 1),
    firstThreeChars: authSecret.slice(0, 3),
    lastThreeChars: authSecret.slice(-3),
    vercelUrl: process.env.VERCEL_URL ?? null,
    vercelEnv: process.env.VERCEL_ENV ?? null,
  });
};
