"use server";

import { GET, POST, PUT } from "@/lib/api";
import type { WorkoutFormDTO, WorkoutPartDTO } from "./types";
import { toHHmm } from "./utils";

/**
 * 空のWorkoutFormDTOを生成
 */
export const buildEmptyDTO = (): WorkoutFormDTO => ({
  id: null,
  performed_date: "",
  started_at: null,
  ended_at: null,
  place: "",
  note: null,
  condition_level: null,
  workout_part: { id: null, name: null, source: null },
  exercises: Array.from({ length: 3 }, () => ({
    id: null,
    name: "",
    workout_part_id: null,
    is_default: 0,
    sets: Array.from({ length: 5 }, (_, i) => ({
      id: null,
      set_number: i + 1,
      weight_kg: "",
      reps: "",
      note: null,
    })),
  })),
});

/**
 * ワークアウト記録を取得
 * @param date - YYYY-MM-DD形式の日付（省略時は今日のJST日付をバックエンドで使用）
 */
export const fetchWorkoutRecord = async (token: string, date?: string): Promise<WorkoutFormDTO> => {
  // dateがある場合はクエリパラメータに含める、ない場合はバックエンドが今日のJST日付を使用
  const url = date ? `/api/v1/workouts/records?date=${date}` : "/api/v1/workouts/records";
  const res = await GET(url, { token: token ?? undefined });

  if (!res.ok) {
    // エラーの場合は空のDTOを返す（日付はバックエンドで設定されるはず）
    return buildEmptyDTO();
  }

  const display = res.data as WorkoutFormDTO;

  // バックエンドから返されたデータをそのまま使用（日付も含まれている）
  if (!display || !display.id) {
    // レコードがない場合もバックエンドから日付とplace、exercisesは返されている
    // 空の3行のエクササイズを作成
    const emptyExercises = Array.from({ length: 3 }, () => ({
      id: null,
      name: "",
      workout_part_id: null,
      is_default: 0 as 0 | 1, // 型を明示的に指定
      sets: Array.from({ length: 5 }, (_, i) => ({
        id: null,
        set_number: i + 1,
        weight_kg: "",
        reps: "",
        note: null,
      })),
    }));

    return {
      id: null,
      performed_date: display?.performed_date || "",
      started_at: null,
      ended_at: null,
      place: display?.place || "",
      note: null,
      condition_level: null,
      workout_part: { id: null, name: null, source: null },
      exercises: display?.exercises?.length > 0 ? display.exercises : emptyExercises,
    };
  }

  return {
    id: display.id,
    performed_date: display.performed_date,
    started_at: toHHmm(display.started_at),
    ended_at: toHHmm(display.ended_at),
    place: display.place ?? "",
    note: display.note ?? null,
    condition_level: display.condition_level ?? null,
    workout_part: {
      id: display.workout_part?.id ?? null,
      name: display.workout_part?.name ?? null,
      source: display.workout_part?.source ?? null,
    },
    exercises: display.exercises.map((ex) => ({
      id: ex.id,
      name: ex.name,
      workout_part_id: ex.workout_part_id,
      is_default: ex.is_default,
      sets: (ex.sets ?? []).map((s) => ({
        id: s.id,
        set_number: s.set_number,
        weight_kg: s.weight_kg,
        reps: s.reps,
        note: s.note ?? null,
      })),
    })),
  };
};

/**
 * ワークアウト部位リストを取得
 */
export const fetchWorkoutParts = async (token: string): Promise<WorkoutPartDTO[]> => {
  const res = await GET<WorkoutPartDTO[]>("/api/v1/workouts/parts", { token: token ?? undefined });

  console.log("[fetchWorkoutParts] Response:", res);
  if (!res.ok || !Array.isArray(res.data)) {
    return [];
  }

  return res.data;
};

/**
 * ワークアウト部位をシード（初期データ作成）
 */
export const seedWorkoutParts = async (token: string): Promise<void> => {
  await POST("/api/v1/workouts/seed", { token: token ?? undefined });
};

// ==================== Server Actions ====================

/**
 * ワークアウト記録を作成
 */
export async function createWorkoutRecord(token: string, body: any) {
  if (!token) {
    return { ok: false, error: "認証エラー" };
  }

  const res = await POST("/api/v1/workouts/records", { body, token: token });
  return { ok: res.ok, error: res.ok ? null : "保存に失敗しました" };
}

/**
 * ワークアウト記録を更新
 */
export async function updateWorkoutRecord(token: string, id: number, body: any) {
  if (!token) {
    return { ok: false, error: "認証エラー" };
  }

  const res = await PUT(`/api/v1/workouts/records/${id}`, { body, token: token });
  return { ok: res.ok, error: res.ok ? null : "更新に失敗しました" };
}
