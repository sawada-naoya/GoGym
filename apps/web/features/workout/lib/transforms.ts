import { formatDate } from "@/utils/time";
import type { WorkoutFormDTO, WorkoutRecordResponseDTO } from "@/types/workout";

/**
 * ISO形式の日時文字列をHH:mm形式に変換
 */
export const toHHmm = (iso: string | null): string | null => {
  if (!iso) return null;
  const d = new Date(iso);
  const hours = d.getHours();
  const minutes = d.getMinutes();
  return `${String(hours).padStart(2, "0")}:${String(minutes).padStart(2, "0")}`;
};

/**
 * 年月日とHH:mm形式の時刻をISO形式に変換
 * @param y - 年
 * @param m - 月
 * @param d - 日
 * @param hm - HH:mm形式の時刻文字列
 * @returns ISO形式の日時文字列 or null
 */
export const toISO = (
  y: number,
  m: number,
  d: number,
  hm?: string | null,
): string | null => {
  if (!hm) return null;
  const [hh, mm] = hm.split(":").map(Number);
  return new Date(y, m - 1, d, hh, mm, 0).toISOString();
};

/**
 * フォームデータを送信用に変換
 * 空のセット（weight_kg と reps が両方とも空）は除外
 *
 * @returns 送信用に変換されたデータ（weight_kg/repsはnumberまたはnull）
 */
export const transformFormDataForSubmit = (
  data: WorkoutFormDTO,
  year: number,
  month: number,
  day: number,
) => ({
  ...data,
  performed_date: formatDate(year, month, day),
  started_at: data.started_at || null,
  ended_at: data.ended_at || null,
  gym_name: data.gym_name?.trim() || null,
  exercises: data.exercises.map((ex) => ({
    ...ex,
    sets: ex.sets
      .map((s) => ({
        ...s,
        id: null, // 新規作成時はIDをnullにする（upsert時にバックエンドで処理）
        weight_kg:
          s.weight_kg === "" || s.weight_kg === null
            ? null
            : Number(s.weight_kg),
        reps: s.reps === "" || s.reps === null ? null : Number(s.reps),
      }))
      .filter((s) => s.weight_kg !== null || s.reps !== null),
  })),
});

/**
 * バックエンドのレスポンスをフロントエンドのフォーム形式に変換
 */
export const convertResponseToFormDTO = (
  response: WorkoutRecordResponseDTO,
): WorkoutFormDTO => {
  // parts配列からexercisesをフラット化
  const exercises = response.parts.flatMap((part) =>
    part.exercises.map((exercise) => ({
      id: exercise.id ?? null,
      name: exercise.name,
      workout_part_id: exercise.workout_part_id ?? null,
      sets: exercise.sets.map((set) => ({
        id: set.id ?? null,
        set_number: set.set_number,
        weight_kg: set.weight_kg ?? "",
        reps: set.reps ?? "",
        note: set.note ?? null,
      })),
    })),
  );

  return {
    id: response.id ?? null,
    performed_date: response.performed_date,
    started_at: response.started_at ?? null,
    ended_at: response.ended_at ?? null,
    gym_id: response.gym_id ?? null,
    gym_name: response.gym_name ?? null,
    note: response.note ?? null,
    condition_level: response.condition_level ?? null,
    exercises: exercises.length > 0 ? exercises : [buildEmptyExercise()],
  };
};

/**
 * 空のWorkoutFormDTOを生成
 */
export const buildEmptyDTO = (): WorkoutFormDTO => ({
  id: null,
  performed_date: "",
  started_at: null,
  ended_at: null,
  gym_id: null,
  gym_name: null,
  note: null,
  condition_level: null,
  exercises: [buildEmptyExercise()],
});

/**
 * 空の種目オブジェクトを生成（内部ヘルパー）
 */
const buildEmptyExercise = () => ({
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
});
