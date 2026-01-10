"use server";

import { revalidatePath } from "next/cache";
import { getServerAccessToken } from "@/features/auth/server";
import type {
  WorkoutFormDTO,
  WorkoutPartDTO,
  ExerciseDTO,
} from "@/types/workout";

const API_BASE = process.env.NEXT_PUBLIC_API_URL;

type ActionResult<T = void> =
  | { success: true; data?: T }
  | { success: false; error: string };

// ==================== Helper ====================

const authorizedFetch = async (
  url: string,
  options?: RequestInit,
): Promise<Response> => {
  const token = await getServerAccessToken();
  if (!token) throw new Error("Unauthorized");
  if (!API_BASE) throw new Error("API base URL not configured");

  return fetch(`${API_BASE}${url}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
      ...options?.headers,
    },
    cache: "no-store",
  });
};

// ==================== Workout Records ====================

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

    const data = await res.json();
    return { success: true, data };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "Unknown error occurred",
    };
  }
};

export const createWorkoutRecord = async (
  body: WorkoutFormDTO,
): Promise<ActionResult> => {
  try {
    const res = await authorizedFetch("/api/v1/workouts/records", {
      method: "POST",
      body: JSON.stringify(body),
    });

    if (!res.ok) {
      return { success: false, error: "Failed to create workout record" };
    }

    revalidatePath("/[locale]/workout", "page");
    return { success: true };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "Unknown error occurred",
    };
  }
};

export const updateWorkoutRecord = async (
  id: string,
  body: WorkoutFormDTO,
): Promise<ActionResult> => {
  try {
    const res = await authorizedFetch(`/api/v1/workouts/records/${id}`, {
      method: "PUT",
      body: JSON.stringify(body),
    });

    if (!res.ok) {
      return { success: false, error: "Failed to update workout record" };
    }

    revalidatePath("/[locale]/workout", "page");
    return { success: true };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "Unknown error occurred",
    };
  }
};

// ==================== Workout Exercises ====================

export const getLastExerciseRecord = async (
  exerciseId: number,
): Promise<ActionResult<ExerciseDTO>> => {
  try {
    const res = await authorizedFetch(
      `/api/v1/workouts/exercises/${exerciseId}/last`,
    );

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

    revalidatePath("/[locale]/workout", "page");
    return { success: true };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "Unknown error occurred",
    };
  }
};

export const deleteWorkoutExercise = async (
  id: number,
): Promise<ActionResult> => {
  try {
    const res = await authorizedFetch(`/api/v1/workouts/exercises/${id}`, {
      method: "DELETE",
    });

    if (!res.ok) {
      return { success: false, error: "Failed to delete exercise" };
    }

    revalidatePath("/[locale]/workout", "page");
    return { success: true };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "Unknown error occurred",
    };
  }
};

// ==================== Workout Parts ====================

export const getWorkoutParts = async (): Promise<
  ActionResult<WorkoutPartDTO[]>
> => {
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

export const seedWorkoutParts = async (): Promise<ActionResult> => {
  try {
    const res = await authorizedFetch("/api/v1/workouts/seed", {
      method: "POST",
    });

    if (!res.ok) {
      return { success: false, error: "Failed to seed workout parts" };
    }

    revalidatePath("/[locale]/workout", "page");
    return { success: true };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : "Unknown error occurred",
    };
  }
};
