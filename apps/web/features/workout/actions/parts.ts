"use server";

import { revalidatePath } from "next/cache";
import { authorizedFetch } from "@/lib/api/client";
import type { ActionResult } from "@/lib/api/types";
import type { WorkoutPartDTO } from "@/types/workout";

/**
 * ワークアウト部位一覧を取得
 */
export const getWorkoutParts = async (): Promise<ActionResult<WorkoutPartDTO[]>> => {
  try {
    const res = await authorizedFetch("/api/v1/workouts/parts");

    if (!res.ok) {
      return { success: false, error: "Failed to fetch workout parts" };
    }

    const data = await res.json();
    return { success: true, data };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "Unknown error occurred",
    };
  }
};

/**
 * ワークアウト部位の初期データを作成（idempotent）
 */
export const seedWorkoutParts = async (): Promise<ActionResult> => {
  try {
    const res = await authorizedFetch("/api/v1/workouts/seed", {
      method: "POST",
    });

    if (!res.ok) {
      return { success: false, error: "Failed to seed workout parts" };
    }

    // 成功したらページを再検証
    revalidatePath("/workout");

    return { success: true };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "Unknown error occurred",
    };
  }
};
