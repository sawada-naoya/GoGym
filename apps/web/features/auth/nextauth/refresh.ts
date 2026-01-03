import type { RefreshResponse } from "../../../types/refresh";

const API_BASE = process.env.NEXT_PUBLIC_API_URL;

/**
 * リフレッシュトークンを使って新しいアクセストークンを取得
 */
export const refreshAccessToken = async (refreshToken: string): Promise<RefreshResponse | null> => {
  try {
    console.log("[DEBUG] refreshAccessToken called");
    console.log("[DEBUG] API_BASE:", API_BASE);
    console.log("[DEBUG] refreshToken exists:", !!refreshToken);

    if (!API_BASE) {
      console.error("[DEBUG] API_BASE is undefined");
      return null;
    }

    const url = `${API_BASE}/api/v1/sessions/refresh`;
    console.log("[DEBUG] Fetching URL:", url);

    const res = await fetch(url, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
      cache: "no-store",
    });

    console.log("[DEBUG] Response status:", res.status);

    if (!res.ok) {
      const errorText = await res.text();
      console.error("[DEBUG] Failed to refresh token:", res.status, errorText);
      return null;
    }

    const data = (await res.json()) as RefreshResponse;
    console.log("[DEBUG] Refresh successful");
    return data;
  } catch (error) {
    console.error("[DEBUG] Error refreshing token:", error);
    if (error instanceof Error) {
      console.error("[DEBUG] Error message:", error.message);
    }
    return null;
  }
};
