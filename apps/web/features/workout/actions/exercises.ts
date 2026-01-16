"use server";

import { revalidatePath } from "next/cache";
import { authorizedFetch } from "@/lib/api/client";
import type { ActionResult } from "@/lib/api/types";
import type { ExerciseDTO } from "@/types/workout";

/**
 * 特定の種目の前回記録を取得
 */
export const getLastExerciseRecord = async (exerciseId: number): Promise<ActionResult<ExerciseDTO>> => {
  try {
    const res = await authorizedFetch(`/api/v1/workouts/exercises/${exerciseId}/last`);

    if (!res.ok || res.status === 204) {
      return { success: false, error: "No previous record found" };
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
 * ワークアウト種目を一括作成/更新
 */
export const upsertWorkoutExercises = async (body: {
  exercises: Array<{
    id?: number | null;
    name: string;
    workout_part_id: number;
  }>;
}): Promise<ActionResult> => {
  try {
    const res = await authorizedFetch("/api/v1/workouts/exercises", {
      method: "POST",
      body: JSON.stringify(body),
    });

    if (!res.ok) {
      return { success: false, error: "Failed to upsert exercises" };
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

/**
 * ワークアウト種目を削除
 */
export const deleteWorkoutExercise = async (id: number): Promise<ActionResult> => {
  try {
    const res = await authorizedFetch(`/api/v1/workouts/exercises/${id}`, {
      method: "DELETE",
    });

    if (!res.ok) {
      return { success: false, error: "Failed to delete exercise" };
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
