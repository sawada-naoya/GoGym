// API base URL
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";

export type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

export interface RequestOptions {
  query?: Record<string, string | number | boolean>;
  headers?: Record<string, string>;
  body?: unknown;
  cache?: RequestCache; // 'no-store' | 'force-cache' など
}

// queryパラメータをURLエンコードして文字列に変換するヘルパー関数
// 例: { search: "gym", page: 2 } → "search=gym&page=2"
const buildQueryParams = (query?: RequestOptions["query"]): string => {
  if (!query) return "";
  const params = new URLSearchParams();
  Object.entries(query).forEach(([key, value]) => {
    params.append(key, String(value));
  });
  return params.toString();
};

const request = async <T = any>(endpoint: string, method: HttpMethod, options: RequestOptions = {}): Promise<Response> => {
  let url = `${API_BASE_URL}${endpoint}`;
  const qs = buildQueryParams(options.query);
  if (qs) url += `?${qs}`;

  const config: RequestInit = {
    method,
    headers: {
      "Content-Type": "application/json",
      ...(options.headers ?? {}),
    },
    cache: options.cache ?? "no-store",
  };

  if (method !== "GET" && options.body) {
    config.body = JSON.stringify(options.body);
  }

  return fetch(url, config);
};

// CRUDエイリアスを大文字関数で用意
export const GET = <T>(endpoint: string, options?: RequestOptions) => request<T>(endpoint, "GET", options);

export const POST = <T>(endpoint: string, body?: unknown, options?: RequestOptions) => request<T>(endpoint, "POST", { ...options, body });

export const PUT = <T>(endpoint: string, body?: unknown, options?: RequestOptions) => request<T>(endpoint, "PUT", { ...options, body });

export const PATCH = <T>(endpoint: string, body?: unknown, options?: RequestOptions) => request<T>(endpoint, "PATCH", { ...options, body });

export const DELETE = <T>(endpoint: string, options?: RequestOptions) => request<T>(endpoint, "DELETE", options);
