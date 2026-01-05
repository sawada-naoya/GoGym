// app/api/_debug/cookies/route.ts
import { cookies } from "next/headers";
import { NextResponse } from "next/server";

export const runtime = "nodejs";
export const dynamic = "force-dynamic";

export const GET = async () => {
  const names = (await cookies()).getAll().map((c) => c.name);
  return NextResponse.json({
    cookieNames: names,
    hasAuthjsSession: names.some((n) => n.includes("authjs") && n.includes("session-token")),
    host: process.env.VERCEL_URL ?? null,
  });
};
