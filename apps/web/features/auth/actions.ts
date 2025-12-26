"use server";

import type { SignUpPayload } from "./schemas/signup";

const API_BASE =
  process.env.API_INTERNAL_URL ?? process.env.NEXT_PUBLIC_API_URL;

type ActionResult<T = void> =
  | { success: true; data?: T }
  | { success: false; error: string };

export const signup = async (payload: SignUpPayload): Promise<ActionResult> => {
  try {
    if (!API_BASE) {
      return { success: false, error: "API base URL not configured" };
    }

    const res = await fetch(`${API_BASE}/api/v1/users`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
      cache: "no-store",
    });

    if (!res.ok) {
      const data = await res.json().catch(() => null);
      return {
        success: false,
        error: data?.message || "このメールアドレスは既に使用されています",
      };
    }

    return { success: true };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "Unknown error occurred",
    };
  }
};
