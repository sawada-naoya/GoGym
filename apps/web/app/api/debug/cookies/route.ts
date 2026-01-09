// src/app/api/_debug/cookies/route.ts
import { cookies, headers } from "next/headers";
import { NextResponse } from "next/server";
import { getToken } from "next-auth/jwt";
import crypto from "crypto";

export const runtime = "nodejs";
export const dynamic = "force-dynamic";

const sha256_12 = (v: string) => crypto.createHash("sha256").update(v).digest("hex").slice(0, 12);
const mask = (v: string, head = 12, tail = 8) => {
  if (!v) return "";
  if (v.length <= head + tail) return `${v.slice(0, 2)}***${v.slice(-2)}`;
  return `${v.slice(0, head)}…${v.slice(-tail)}(len=${v.length})`;
};

export const GET = async () => {
  const h = await headers();
  const host = h.get("host");
  const forwardedHost = h.get("x-forwarded-host");
  const proto = h.get("x-forwarded-proto");

  const store = await cookies();
  const all = store.getAll();

  const cookieNames = all.map((c) => c.name);
  const hasAuthjsSession = cookieNames.some((n) => n.includes("authjs") && n.includes("session-token"));
  const hasNextAuthSession = cookieNames.some((n) => n.includes("next-auth") && n.includes("session-token"));

  // cookieHeader を作って getToken に渡す
  const cookieHeader = all.map((c) => `${c.name}=${c.value}`).join("; ");

  // セッションcookie候補の長さだけ（値は返さない）
  const sessionCandidates = ["__Secure-authjs.session-token", "__Host-authjs.session-token", "authjs.session-token", "__Secure-next-auth.session-token", "next-auth.session-token"];

  const sessionCookieInfo = sessionCandidates
    .map((name) => {
      const v = store.get(name)?.value;
      return v ? { name, len: v.length, hash12: sha256_12(v) } : null;
    })
    .filter(Boolean);

  const authSecret = process.env.AUTH_SECRET ?? "";
  const nextAuthSecret = process.env.NEXTAUTH_SECRET ?? "";
  const effectiveSecret = authSecret || nextAuthSecret;

  // token decode（失敗しても throw しない）
  let token: any = null;
  let decodeErr: any = null;

  try {
    if (effectiveSecret) {
      token = await getToken({ req: { headers: { cookie: cookieHeader } } as any, secret: effectiveSecret });
    }
  } catch (e: any) {
    decodeErr = { name: e?.name, message: e?.message };
  }

  // ✅ 詳細はレスポンスに載せず、ログだけに出す
  console.error("[DEBUG_COOKIES]", {
    request: {
      host,
      forwardedHost,
      proto,
      vercelEnv: process.env.VERCEL_ENV ?? null,
      vercelUrl: process.env.VERCEL_URL ?? null,
      nodeEnv: process.env.NODE_ENV ?? null,
      userAgent: mask(h.get("user-agent") ?? "", 40, 0),
      referer: mask(h.get("referer") ?? "", 60, 0),
    },
    env: {
      hasAUTH_SECRET: Boolean(authSecret),
      hasNEXTAUTH_SECRET: Boolean(nextAuthSecret),
      authSecretFp12: authSecret ? sha256_12(authSecret) : null, // secretそのものは出さない
      nextAuthSecretFp12: nextAuthSecret ? sha256_12(nextAuthSecret) : null,
      authUrl: process.env.AUTH_URL ?? process.env.NEXTAUTH_URL ?? null,
    },
    cookies: {
      cookieNames,
      hasAuthjsSession,
      hasNextAuthSession,
      cookieHeaderLen: cookieHeader.length,
      sessionCookieInfo,
      // 全cookieの値は危険なので “長さ” と “hash” だけ
      cookieStats: all.map((c) => ({ name: c.name, len: c.value.length, hash12: sha256_12(c.value) })),
    },
    token: token
      ? {
          tokenIsNull: false,
          keys: Object.keys(token),
          sub: token.sub ?? null,
          email: token.email ?? null,
          hasAccessToken: Boolean(token.accessToken),
          hasRefreshToken: Boolean(token.refreshToken),
          accessTokenMasked: token.accessToken ? mask(String(token.accessToken), 16, 10) : null,
          refreshTokenMasked: token.refreshToken ? mask(String(token.refreshToken), 16, 10) : null,
          expiresAt: token.expiresAt ?? null,
          iat: token.iat ?? null,
          exp: token.exp ?? null,
        }
      : {
          tokenIsNull: true,
          decodeErr,
        },
    diagnosis: (() => {
      if (!hasAuthjsSession && !hasNextAuthSession) return "NO_SESSION_COOKIE_SENT";
      if (!effectiveSecret) return "NO_SECRET_AT_RUNTIME";
      if (!token) return decodeErr ? "DECRYPT_THROWN_SECRET_OR_FORMAT" : "DECRYPT_FAILED_SECRET_MISMATCH_OR_COOKIE_TRUNCATED";
      if (!token.accessToken) return "TOKEN_DECODED_BUT_NO_ACCESS_TOKEN_IN_JWT_CALLBACK";
      return "OK";
    })(),
  });

  // ✅ レスポンスは最小限に留める（公開されても致命傷にならない範囲）
  return NextResponse.json({
    cookieNames,
    hasAuthjsSession,
    host: process.env.VERCEL_URL ?? null,
    tokenDecoded: Boolean(token),
    hasAccessToken: Boolean(token?.accessToken),
    // “危険な値” は返さない
  });
};
