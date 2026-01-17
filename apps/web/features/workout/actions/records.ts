"use server";

import { revalidatePath } from "next/cache";
import { authorizedFetch } from "@/lib/api/client";
import type { ActionResult } from "@/lib/api/types";
import type { WorkoutFormDTO, WorkoutRecordResponseDTO } from "@/types/workout";
import { convertResponseToFormDTO } from "@/features/workout/lib/transforms";

type WorkoutSubmitDTO = Omit<WorkoutFormDTO, "exercises"> & {
  exercises: Array<{
    id?: number | null;
    name: string;
    workout_part_id: number | null;
    sets: Array<{
      id?: number | null;
      set_number: number;
      weight_kg: number | null;
      reps: number | null;
      note: string | null;
    }>;
  }>;
};

/**
 * ワークアウトレコードを取得
 */
export const getWorkoutRecords = async (params?: {
  date?: string;
  partId?: number;
}): Promise<ActionResult<WorkoutFormDTO>> => {
  try {
    const searchParams = new URLSearchParams();
    if (params?.date) searchParams.set("date", params.date);
    if (params?.partId) searchParams.set("part_id", String(params.partId));

    const queryString = searchParams.toString()
      ? `?${searchParams.toString()}`
      : "";

    const res = await authorizedFetch(`/api/v1/workouts/records${queryString}`);

    if (!res.ok) {
      return { success: false, error: "Failed to fetch workout records" };
    }

    const response: WorkoutRecordResponseDTO = await res.json();
    const data = convertResponseToFormDTO(response);
    return { success: true, data };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "Unknown error occurred",
    };
  }
};

/**
 * ワークアウトレコードを新規作成
 */
export const createWorkoutRecord = async (
  body: WorkoutSubmitDTO,
): Promise<ActionResult> => {
  try {
    const res = await authorizedFetch("/api/v1/workouts/records", {
      method: "POST",
      body: JSON.stringify(body),
    });

    if (!res.ok) {
      return { success: false, error: "Failed to create workout record" };
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
 * ワークアウトレコードを更新
 */
export const updateWorkoutRecord = async (
  id: string,
  body: WorkoutSubmitDTO,
): Promise<ActionResult> => {
  try {
    const res = await authorizedFetch(`/api/v1/workouts/records/${id}`, {
      method: "PUT",
      body: JSON.stringify(body),
    });

    if (!res.ok) {
      return { success: false, error: "Failed to update workout record" };
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
