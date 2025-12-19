import { NextResponse } from "next/server";
import { getServerAccessToken } from "@/lib/auth-helpers";

const API_BASE =
  process.env.API_INTERNAL_URL ?? process.env.NEXT_PUBLIC_API_URL;

type ProxyOptions = {
  method?: string;
  body?: any;
  headers?: Record<string, string>;
};

/**
 * 認証が必要なエンドポイント用のプロキシ関数
 * @param req - Next.js Request object
 * @param path - Go API のパス（例: "/api/v1/workouts/records"）
 * @param options - メソッド、body、追加ヘッダー（bodyはオブジェクトの場合自動でstringifyされる）
 */
export const proxyToGo = async (
  req: Request,
  path: string,
  options?: ProxyOptions,
) => {
  const token = await getServerAccessToken();
  if (!token) {
    return errorResponse("unauthorized", 401);
  }
  return proxyRequest(req, path, options, token);
};

/**
 * 認証不要のエンドポイント用のプロキシ関数（signup等）
 * @param req - Next.js Request object
 * @param path - Go API のパス（例: "/api/v1/users"）
 * @param options - メソッド、body、追加ヘッダー（bodyはオブジェクトの場合自動でstringifyされる）
 */
export const proxyToGoPublic = async (
  req: Request,
  path: string,
  options?: ProxyOptions,
) => {
  return proxyRequest(req, path, options);
};

/**
 * 共通のプロキシ処理
 */
const proxyRequest = async (
  req: Request,
  path: string,
  options?: ProxyOptions,
  token?: string,
) => {
  if (!API_BASE) {
    return errorResponse("API base url is not configured", 500);
  }

  const url = new URL(req.url);
  const target = `${API_BASE}${path}${url.search}`;

  const headers: Record<string, string> = {
    ...(options?.headers ?? {}),
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  let body: string | undefined;
  if (options?.body !== undefined) {
    body =
      typeof options.body === "string"
        ? options.body
        : JSON.stringify(options.body);
    headers["content-type"] = "application/json";
  }

  const res = await fetch(target, {
    method: options?.method ?? "GET",
    headers,
    body,
    cache: "no-store",
  });

  const data = await res.json().catch(() => null);
  return NextResponse.json(data, { status: res.status });
};

/**
 * エラーレスポンス用ヘルパー
 */
export const errorResponse = (message: string, status: number) => {
  return NextResponse.json({ message }, { status });
};
