type ApiSuccess<T> = {
  ok: true;
  status: number;
  data: T;
  headers: Headers;
};

type ApiFailure<E = unknown> = {
  ok: false;
  status: number;
  error?: E;
  headers: Headers;
};

type ApiResponse<T, E = unknown> = ApiSuccess<T> | ApiFailure<E>;

type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

export type Query = Record<string, string | number | boolean | (string | number | boolean)[] | null | undefined>;

type RequestOptions = {
  query?: Query;
  body?: unknown;
  headers?: Record<string, string>;
};

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL;

const isJson = (h: Headers) => (h.get("content-type") || "").toLowerCase().includes("application/json");

const appendQuery = (endpoint: string, q?: Query) => {
  if (!q) return endpoint;
  const u = new URL(endpoint);
  u.searchParams.sort();
  for (const [k, v] of Object.entries(q)) {
    if (v === undefined || v === null) continue;
    if (Array.isArray(v)) v.forEach((item) => u.searchParams.append(k, String(item)));
    else u.searchParams.set(k, String(v));
  }
  return u.pathname + (u.search ? `?${u.searchParams.toString()}` : "");
};

const _request = async <T, E = { message?: string }>(method: HttpMethod, endpoint: string, options: RequestOptions = {}): Promise<ApiResponse<T, E>> => {
  const path = appendQuery(endpoint, options.query);
  const url = API_BASE_URL + path;

  const headers: Record<string, string> = { ...(options.headers ?? {}) };
  const hasBody = method !== "GET" && options.body !== undefined;
  if (hasBody && !headers["Content-Type"]) headers["Content-Type"] = "application/json";

  const res = await fetch(url, {
    method,
    headers,
    body: hasBody ? JSON.stringify(options.body) : undefined,
    credentials: "include",
  });
  // JSONの場合はレスポンスをパース
  const parse = async <X>(): Promise<X | undefined> => {
    if (res.status === 204) return undefined;
    if (!isJson(res.headers)) return undefined;
    try {
      return (await res.json()) as X;
    } catch {
      return undefined;
    }
  };
  if (res.ok) {
    const data = (await parse<T>()) as T;
    return { ok: true, status: res.status, data, headers: res.headers };
  } else {
    const error = await parse<E>();
    return { ok: false, status: res.status, error, headers: res.headers };
  }
};

// CRUDエイリアス
export const GET = <T>(endpoint: string, options?: RequestOptions) => _request<T>("GET", endpoint, options);

export const POST = <T>(endpoint: string, options?: RequestOptions) => _request<T>("POST", endpoint, options);

export const PUT = <T>(endpoint: string, options?: RequestOptions) => _request<T>("PUT", endpoint, options);

export const PATCH = <T>(endpoint: string, options?: RequestOptions) => _request<T>("PATCH", endpoint, options);

export const DELETE_ = <T>(endpoint: string, options?: RequestOptions) => _request<T>("DELETE", endpoint, options);
