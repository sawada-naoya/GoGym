import type { RefreshResponse } from "../../../types/refresh";

const API_BASE =
  process.env.API_INTERNAL_URL ?? process.env.NEXT_PUBLIC_API_URL;

/**
 * リフレッシュトークンを使って新しいアクセストークンを取得
 */
export const refreshAccessToken = async (
  refreshToken: string,
): Promise<RefreshResponse | null> => {
  try {
    const res = await fetch(`${API_BASE}/api/v1/sessions/refresh`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
      cache: "no-store",
    });

    if (!res.ok) {
      console.error("Failed to refresh token:", res.status);
      return null;
    }

    const data = (await res.json()) as RefreshResponse;
    return data;
  } catch (error) {
    console.error("Error refreshing token:", error);
    return null;
  }
};
