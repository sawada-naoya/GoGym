import { getServerAccessToken } from "@/features/auth/server";

const API_BASE = process.env.NEXT_PUBLIC_API_URL;

/**
 * APIベースURLと認証チェック
 */
const validateApiConfig = () => {
  if (!API_BASE) {
    throw new Error("API base URL not configured");
  }
  return API_BASE;
};

/**
 * 認証トークン付きfetch
 * Server Actions/Server Componentsでのみ使用可能
 *
 * @param url - APIエンドポイント（/api/v1/...）
 * @param options - fetch options
 * @throws {Error} トークンがない、またはAPI_BASEが未設定の場合
 */
export const authorizedFetch = async (url: string, options?: RequestInit): Promise<Response> => {
  const token = await getServerAccessToken();
  if (!token) {
    throw new Error("Unauthorized: No access token available");
  }

  const apiBase = validateApiConfig();

  return fetch(`${apiBase}${url}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
      ...options?.headers,
    },
    cache: "no-store",
  });
};

/**
 * 認証不要のfetch（signup等）
 * Server Actions/Server Componentsでのみ使用可能
 *
 * @param url - APIエンドポイント（/api/v1/...）
 * @param options - fetch options
 * @throws {Error} API_BASEが未設定の場合
 */
export const apiFetch = async (url: string, options?: RequestInit): Promise<Response> => {
  const apiBase = validateApiConfig();

  return fetch(`${apiBase}${url}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
    cache: "no-store",
  });
};
