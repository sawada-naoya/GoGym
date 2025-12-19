import type {
  WorkoutFormDTO,
  WorkoutPartDTO,
  ExerciseDTO,
} from "@/app/workout/_lib/types";

/**
 * ワークアウト部位をシード（初期データ作成）
 */
export const seedWorkoutParts = async (): Promise<void> => {
  const res = await fetch("/api/workouts/seed", {
    method: "POST",
    cache: "no-store",
  });
  if (!res.ok) throw new Error("failed to seed workout parts");
};

/**
 * ワークアウト部位リストを取得
 */
export const getWorkoutParts = async (): Promise<WorkoutPartDTO[]> => {
  const res = await fetch("/api/workouts/parts", { cache: "no-store" });
  if (!res.ok) return [];
  return res.json();
};

/**
 * ワークアウト記録を取得
 */
export const getWorkoutRecord = async (params?: {
  date?: string;
  partId?: number | null;
}): Promise<WorkoutFormDTO> => {
  const q = new URLSearchParams();
  if (params?.date) q.set("date", params.date);
  if (params?.partId) q.set("part_id", String(params.partId));

  const url = q.toString()
    ? `/api/workouts/records?${q.toString()}`
    : "/api/workouts/records";
  const res = await fetch(url, { cache: "no-store" });
  if (!res.ok) throw new Error("failed to fetch workout record");
  return res.json();
};

/**
 * ワークアウト記録を作成
 */
export const createWorkoutRecord = async (
  body: any,
): Promise<{ ok: boolean; error: string | null }> => {
  const res = await fetch("/api/workouts/records", {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify(body),
    cache: "no-store",
  });
  return { ok: res.ok, error: res.ok ? null : "保存に失敗しました" };
};

/**
 * ワークアウト記録を更新
 */
export const updateWorkoutRecord = async (
  id: number,
  body: any,
): Promise<{ ok: boolean; error: string | null }> => {
  const res = await fetch(`/api/workouts/records/${id}`, {
    method: "PUT",
    headers: { "content-type": "application/json" },
    body: JSON.stringify(body),
    cache: "no-store",
  });
  return { ok: res.ok, error: res.ok ? null : "更新に失敗しました" };
};

/**
 * 前回のエクササイズ記録を取得
 */
export const getLastExerciseRecord = async (
  exerciseID: number,
): Promise<ExerciseDTO | null> => {
  if (!exerciseID) {
    return null;
  }

  const res = await fetch(
    `/api/workouts/exercises?id=${exerciseID}&action=last`,
    {
      cache: "no-store",
    },
  );

  if (!res.ok || res.status === 204) {
    return null;
  }

  const data = await res.json();
  return data || null;
};

/**
 * ワークアウト種目を一括作成/更新（Upsert）
 */
export const upsertWorkoutExercises = async (
  exercises: Array<{
    id?: number;
    name: string;
    workout_part_id: number | null;
  }>,
): Promise<{ ok: boolean; error: string | null; data?: any }> => {
  const res = await fetch("/api/workouts/exercises", {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({ exercises }),
    cache: "no-store",
  });

  const data = res.ok ? await res.json().catch(() => null) : null;
  return {
    ok: res.ok,
    error: res.ok ? null : "種目の登録に失敗しました",
    data,
  };
};

/**
 * ワークアウト種目を削除
 */
export const deleteWorkoutExercise = async (
  exerciseID: number,
): Promise<{ ok: boolean; error: string | null }> => {
  const res = await fetch(`/api/workouts/exercises?id=${exerciseID}`, {
    method: "DELETE",
    cache: "no-store",
  });

  return { ok: res.ok, error: res.ok ? null : "種目の削除に失敗しました" };
};
