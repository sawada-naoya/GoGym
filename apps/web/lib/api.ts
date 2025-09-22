// API base URL
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL;

export type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
export type Query = Record<string, string | number | boolean | undefined>;

export type RequestOptions = Omit<RequestInit, "method" | "body"> & {
  body?: unknown;
  query?: Query;
};

export type ApiResponse<T> = {
  ok: boolean;
  status: number;
  data: T | null;
};

const buildQueryParams = (query?: Query): string => {
  if (!query) return "";
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(query)) {
    if (value === undefined) continue;
    params.append(key, String(value));
  }
  return params.toString();
};

const parseJsonIfAny = async (res: Response): Promise<unknown | undefined> => {
  const ct = res.headers.get("Content-Type") || "";
  if (!ct.includes("application/json")) return undefined;
  try {
    const text = await res.text();
    if (!text) return undefined;
    return JSON.parse(text);
  } catch {
    return undefined;
  }
};

const request = async <T>(endpoint: string, method: HttpMethod, options: RequestOptions = {}): Promise<ApiResponse<T>> => {
  let url = API_BASE_URL ? `${API_BASE_URL}${endpoint}` : endpoint;
  const qs = buildQueryParams(options.query);
  if (qs) url += (url.includes("?") ? "&" : "?") + qs;

  const res = await fetch(url, {
    ...options,
    method,
    headers: {
      "Content-Type": "application/json",
      ...(options.headers ?? {}),
    },
    body: method !== "GET" && options.body !== undefined ? JSON.stringify(options.body) : undefined,
  });

  let data: T | null = null;
  try {
    data = await res.json();
  } catch {
    data = null;
  }

  return { ok: res.ok, status: res.status, data };
};

// CRUDエイリアス
export const GET = <T>(endpoint: string, options?: RequestOptions) => request<T>(endpoint, "GET", options);

export const POST = <T>(endpoint: string, options?: RequestOptions) => request<T>(endpoint, "POST", options);

export const PUT = <T>(endpoint: string, options?: RequestOptions) => request<T>(endpoint, "PUT", options);

export const PATCH = <T>(endpoint: string, options?: RequestOptions) => request<T>(endpoint, "PATCH", options);

export const DELETE_ = <T>(endpoint: string, options?: RequestOptions) => request<T>(endpoint, "DELETE", options);
