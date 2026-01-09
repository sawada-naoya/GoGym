// src/app/api/_debug/cookies/route.ts
import { cookies, headers } from "next/headers";
import { NextResponse } from "next/server";
import { getToken } from "next-auth/jwt";

export const runtime = "nodejs";
export const dynamic = "force-dynamic";

export const GET = async () => {
  const h = await headers();
  const store = await cookies();
  const all = store.getAll();

  const cookieNames = all.map((c) => c.name);
  const hasAuthjsSession = cookieNames.some((n) => n.includes("authjs") && n.includes("session-token"));

  // session cookie の長さ（破損/切断の疑いを見たい）
  const sessionCookie = store.get("__Secure-authjs.session-token") ?? store.get("__Host-authjs.session-token") ?? store.get("authjs.session-token");

  const sessionCookieLen = sessionCookie?.value.length ?? 0;

  // 復号に使う secret の存在確認（値は返さない）
  const authSecret = process.env.AUTH_SECRET ?? "";
  const nextAuthSecret = process.env.NEXTAUTH_SECRET ?? "";
  const effectiveSecret = authSecret || nextAuthSecret;

  const effectiveSecretSource = authSecret ? "AUTH_SECRET" : nextAuthSecret ? "NEXTAUTH_SECRET" : "NONE";

  // AUTH_SECRETの検証情報（実際の値は絶対に返さない）
  const secretInfo = effectiveSecret
    ? {
        length: effectiveSecret.length,
        firstCharCode: effectiveSecret.charCodeAt(0),
        lastCharCode: effectiveSecret.charCodeAt(effectiveSecret.length - 1),
        hasLeadingQuote: effectiveSecret.startsWith('"') || effectiveSecret.startsWith("'"),
        hasTrailingQuote: effectiveSecret.endsWith('"') || effectiveSecret.endsWith("'"),
        hasLeadingSpace: effectiveSecret[0] === " ",
        hasTrailingSpace: effectiveSecret[effectiveSecret.length - 1] === " ",
        hasNewline: effectiveSecret.includes("\n") || effectiveSecret.includes("\r"),
      }
    : null;

  const cookieHeader = all.map((c) => `${c.name}=${c.value}`).join("; ");

  let token: any = null;
  let decodeError: string | null = null;

  try {
    if (effectiveSecret) {
      token = await getToken({ req: { headers: { cookie: cookieHeader } } as any, secret: effectiveSecret });
    } else {
      decodeError = "NO_SECRET";
    }
  } catch (e: any) {
    decodeError = `${e?.name ?? "Error"}: ${e?.message ?? ""}`.slice(0, 160);
  }

  const diagnosis = (() => {
    if (!hasAuthjsSession) return "NO_SESSION_COOKIE_SENT";
    if (!effectiveSecret) return "NO_SECRET_AT_RUNTIME";
    if (!token) return "DECRYPT_FAILED_SECRET_MISMATCH_OR_COOKIE_TRUNCATED";
    if (!token.accessToken) return "TOKEN_DECODED_BUT_NO_ACCESS_TOKEN_IN_JWT_CALLBACK";
    return "OK";
  })();

  return NextResponse.json({
    cookieNames,
    hasAuthjsSession,
    host: process.env.VERCEL_URL ?? null,

    // 追加：原因特定用
    vercelEnv: process.env.VERCEL_ENV ?? null,
    reqHost: h.get("host"),
    forwardedHost: h.get("x-forwarded-host"),
    proto: h.get("x-forwarded-proto"),

    hasAUTH_SECRET: Boolean(authSecret),
    authSecret: authSecret,
    hasNEXTAUTH_SECRET: Boolean(nextAuthSecret),
    effectiveSecretSource,
    secretInfo,

    sessionCookieLen,
    cookieHeaderLen: cookieHeader.length,

    tokenDecoded: Boolean(token),
    hasAccessToken: Boolean(token?.accessToken),
    decodeError,
    diagnosis,
  });
};
