// API base URL
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL;

export type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
export type Query = Record<string, string | number | boolean | undefined>;

export interface RequestOptions {
  query?: Query; // クエリパラメータ
  headers?: Record<string, string>;
  body?: unknown;
  cache?: RequestCache; // 'no-store' | 'force-cache' など
}

// queryパラメータをURLエンコードして文字列に変換するヘルパー関数
// 例: { search: "gym", page: 2 } → "search=gym&page=2"
const buildQueryParams = (query?: Query): string => {
  if (!query) return "";
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(query)) {
    if (value === undefined) continue;
    params.append(key, String(value));
  }
  return params.toString();
};

const request = async (endpoint: string, method: HttpMethod, options: RequestOptions = {}): Promise<Response> => {
  let url = `${API_BASE_URL}${endpoint}`;
  const queryParams = buildQueryParams(options.query);
  if (queryParams) url += (url.includes("?") ? "&" : "?") + queryParams;
  const init: RequestInit & { next?: NextFetchRequestConfig } = {
    method,
    headers: {
      "Content-Type": "application/json",
      ...(options.headers ?? {}),
    },
    cache: options.cache ?? "no-store",
    ...(method !== "GET" && options.body ? { body: JSON.stringify(options.body) } : {}),
  };
  const res = await fetch(url, init);
  return res;
};

// CRUDエイリアスを大文字関数で用意
export const GET = async <T>(endpoint: string, options?: RequestOptions): Promise<T> => {
  const response = await request(endpoint, "GET", options);
  return response.json();
};

export const POST = async <T>(endpoint: string, body?: unknown, options?: RequestOptions): Promise<T> => {
  const response = await request(endpoint, "POST", { ...options, body });
  return response.json();
};

export const PUT = async <T>(endpoint: string, body?: unknown, options?: RequestOptions): Promise<T> => {
  const response = await request(endpoint, "PUT", { ...options, body });
  return response.json();
};

export const PATCH = async <T>(endpoint: string, body?: unknown, options?: RequestOptions): Promise<T> => {
  const response = await request(endpoint, "PATCH", { ...options, body });
  return response.json();
};

export const DELETE = async <T>(endpoint: string, options?: RequestOptions): Promise<T> => {
  const response = await request(endpoint, "DELETE", options);
  return response.json();
};
