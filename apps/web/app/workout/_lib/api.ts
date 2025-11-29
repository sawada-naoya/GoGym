"use server";

import { GET, POST, PUT, DELETE_ } from "@/lib/api";
import type { WorkoutFormDTO, WorkoutPartDTO } from "./types";
import { buildEmptyDTO } from "./types";
import { ensureFiveSets } from "./utils";

/**
 * ワークアウト記録を取得
 * @param date - YYYY-MM-DD形式の日付（省略時は今日のJST日付をバックエンドで使用）
 * @param partID - 部位ID（指定時はその部位のみのレコードを取得）
 */
export const fetchWorkoutRecord = async (token: string, date?: string, partID?: number | null): Promise<WorkoutFormDTO> => {
  // クエリパラメータを組み立て
  const params = new URLSearchParams();
  if (date) params.append("date", date);
  if (partID) params.append("part_id", partID.toString());

  const queryString = params.toString();
  const url = queryString ? `/api/v1/workouts/records?${queryString}` : "/api/v1/workouts/records";
  const res = await GET(url, { token: token ?? undefined });

  if (!res.ok) {
    // エラーの場合は空のDTOを返す（日付はバックエンドで設定されるはず）
    return buildEmptyDTO();
  }

  const display = res.data as WorkoutFormDTO;

  // バックエンドから返されたデータをそのまま使用（日付も含まれている）
  if (!display || !display.id) {
    // レコードがない場合もバックエンドから日付とplace、exercisesは返されている
    // 空の1行のエクササイズを作成
    const emptyExercises = Array.from({ length: 1 }, () => ({
      id: null,
      name: "",
      workout_part_id: null,
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
    started_at: display.started_at ?? null,
    ended_at: display.ended_at ?? null,
    place: display.place ?? "",
    note: display.note ?? null,
    condition_level: display.condition_level ?? null,
    workout_part: {
      id: display.workout_part?.id ?? null,
      name: display.workout_part?.name ?? null,
      source: display.workout_part?.source ?? null,
    },
    exercises: display.exercises.map((ex) => {
      const exercise = {
        id: ex.id,
        name: ex.name,
        workout_part_id: ex.workout_part_id,
        sets: (ex.sets ?? []).map((s) => ({
          id: s.id,
          set_number: s.set_number,
          weight_kg: s.weight_kg,
          reps: s.reps,
          note: s.note ?? null,
        })),
      };
      // 5セットに揃える
      return ensureFiveSets(exercise);
    }),
  };
};

/**
 * ワークアウト部位リストを取得
 */
export const fetchWorkoutParts = async (token: string): Promise<WorkoutPartDTO[]> => {
  const res = await GET<WorkoutPartDTO[]>("/api/v1/workouts/parts", { token: token ?? undefined });

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

// ==================== Workout Exercises CRUD ====================

/**
 * 種目を一括作成/更新（Upsert）
 */
export async function upsertWorkoutExercises(token: string, exercises: Array<{ id?: number; name: string; workout_part_id: number | null }>) {
  if (!token) {
    return { ok: false, error: "認証エラー" };
  }

  const res = await POST("/api/v1/workouts/exercises/bulk", { body: { exercises }, token });
  return { ok: res.ok, error: res.ok ? null : "種目の登録に失敗しました", data: res.ok ? res.data : null };
}

/**
 * 種目を削除
 */
export async function deleteWorkoutExercise(token: string, exerciseID: number) {
  if (!token) {
    return { ok: false, error: "認証エラー" };
  }

  const res = await DELETE_(`/api/v1/workouts/exercises/${exerciseID}`, { token });
  return { ok: res.ok, error: res.ok ? null : "種目の削除に失敗しました" };
}
