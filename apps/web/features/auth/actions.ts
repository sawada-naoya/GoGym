"use server";

import { apiFetch } from "@/lib/api/client";
import type { ActionResult } from "@/lib/api/types";
import type { SignUpPayload } from "./schemas";

export const signup = async (payload: SignUpPayload): Promise<ActionResult> => {
  try {
    const res = await apiFetch("/api/v1/users", {
      method: "POST",
      body: JSON.stringify(payload),
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
