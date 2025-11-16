"use server";

import { GET, POST, PUT } from "@/lib/api";
import { getToken } from "next-auth/jwt";
import { cookies } from "next/headers";
import type { WorkoutFormDTO, WorkoutPartDTO } from "./types";
import { toHHmm } from "./utils";

/**
 * SSR用: JWTトークンからアクセストークンを取得
 */
async function getAccessToken(): Promise<string | null> {
  const cookieStore = cookies();
  const token = await getToken({
    req: {
      headers: {
        cookie: cookieStore.toString(),
      },
    } as any,
    secret: process.env.NEXTAUTH_SECRET,
  });
  return (token as any)?.accessToken ?? null;
}

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
  exercises: [
    {
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
    },
  ],
});

/**
 * ワークアウト記録を取得
 * @param date - YYYY-MM-DD形式の日付（省略時は今日のJST日付をバックエンドで使用）
 */
export const fetchWorkoutRecord = async (date?: string): Promise<WorkoutFormDTO> => {
  const token = await getAccessToken();

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
    return {
      ...buildEmptyDTO(),
      performed_date: display.performed_date, // バックエンドが設定した日付を使用
      place: display.place || "",
      exercises: display.exercises || buildEmptyDTO().exercises,
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
export const fetchWorkoutParts = async (): Promise<WorkoutPartDTO[]> => {
  const token = await getAccessToken();
  const res = await GET("/api/v1/workouts/parts", { token: token ?? undefined });

  if (!res.ok || !Array.isArray(res.data)) {
    return [];
  }

  return res.data as WorkoutPartDTO[];
};

// ==================== Server Actions ====================

/**
 * ワークアウト記録を作成
 */
export async function createWorkoutRecord(body: any) {
  const accessToken = await getAccessToken();

  if (!accessToken) {
    return { ok: false, error: "認証エラー" };
  }

  const res = await POST("/api/v1/workouts/records", { body, token: accessToken });
  return { ok: res.ok, error: res.ok ? null : "保存に失敗しました" };
}

/**
 * ワークアウト記録を更新
 */
export async function updateWorkoutRecord(id: number, body: any) {
  const accessToken = await getAccessToken();

  if (!accessToken) {
    return { ok: false, error: "認証エラー" };
  }

  const res = await PUT(`/api/v1/workouts/records/${id}`, { body, token: accessToken });
  return { ok: res.ok, error: res.ok ? null : "更新に失敗しました" };
}
