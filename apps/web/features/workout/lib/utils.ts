import type { WorkoutFormDTO, WorkoutPartDTO } from "@/types/workout";
import { formatDate, formatTimeFromDate, parseDateString } from "@/utils/time";

export type ExerciseRow = WorkoutFormDTO["exercises"][number];

// ==================== i18n Helpers ====================

/**
 * Get the localized name for a workout part based on current locale
 * @param part - WorkoutPartDTO with translations
 * @param locale - Target locale (defaults to browser locale or "ja")
 * @returns Localized name string
 */
export const getLocalizedPartName = (
  part: WorkoutPartDTO,
  locale?: string,
): string => {
  // Use provided locale or try to detect from browser
  const targetLocale =
    locale || (typeof navigator !== "undefined" ? navigator.language : "ja");

  // Extract language code (e.g., "en-US" -> "en", "ja-JP" -> "ja")
  const languageCode = targetLocale.split("-")[0];

  // Try to find exact match first
  const exactMatch = part.translations.find((t) => t.locale === targetLocale);
  if (exactMatch) return exactMatch.name;

  // Try to find language code match (e.g., "en" for "en-US")
  const languageMatch = part.translations.find(
    (t) => t.locale.split("-")[0] === languageCode,
  );
  if (languageMatch) return languageMatch.name;

  // Fallback to Japanese
  const jaMatch = part.translations.find((t) => t.locale === "ja");
  if (jaMatch) return jaMatch.name;

  // Fallback to first available translation
  return part.translations[0]?.name || part.key;
};

// ==================== Date & Time Formatters ====================

/**
 * ISO形式の日時文字列をHH:mm形式に変換
 */
export const toHHmm = (iso: string | null): string | null => {
  if (!iso) return null;
  const d = new Date(iso);
  return formatTimeFromDate(d);
};

/**
 * YYYY-MM-DD形式の日付文字列から年月日を抽出
 */
export const extractDateParts = (dateStr: string | null | undefined) => {
  if (!dateStr || dateStr.length < 10) {
    // フォールバック: 今日のJST日付を使用
    const now = new Date();
    return {
      year: now.getFullYear(),
      month: now.getMonth() + 1,
      day: now.getDate(),
    };
  }

  return parseDateString(dateStr);
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
 * 年月日から YYYY-MM-DD 形式の日付文字列を生成
 * @deprecated lib/time.ts の formatDate を使用してください
 */
export const formatDateYMD = (
  year: number,
  month: number,
  day: number,
): string => {
  return formatDate(year, month, day);
};

// ==================== Exercise Helpers ====================

/**
 * Exerciseのセット数を5つに揃える
 */
export const ensureFiveSets = (row: ExerciseRow): ExerciseRow => {
  const sets = row.sets ?? [];
  if (sets.length >= 5) return row;

  const baseLen = sets.length;
  const add: ExerciseRow["sets"] = Array.from(
    { length: 5 - baseLen },
    (_, i) => ({
      set_number: baseLen + i + 1,
      weight_kg: 0,
      reps: 0,
      note: null,
    }),
  );

  return { ...row, sets: [...sets, ...add] };
};

/**
 * 新しい空のExercise行を作成
 * @param setCount - 初期セット数（デフォルト: 5）
 */
export const createEmptyExerciseRow = (setCount: number = 5): ExerciseRow => ({
  id: null,
  name: "",
  workout_part_id: null,
  sets: Array.from({ length: setCount }, (_, i) => ({
    set_number: i + 1,
    weight_kg: 0,
    reps: 0,
    note: null,
  })),
});

/**
 * Exercise行配列のセル値を更新
 */
export const updateExerciseCell = (
  rows: ExerciseRow[],
  rowIndex: number,
  setIndex: number,
  key: "weight_kg" | "reps",
  value: string,
): ExerciseRow[] => {
  const next = structuredClone(rows);
  (next[rowIndex].sets[setIndex] as any)[key] = value;
  return next;
};

/**
 * Exerciseのメモを更新（最初のセットに保存）
 */
export const updateExerciseNote = (
  rows: ExerciseRow[],
  rowIndex: number,
  note: string,
): ExerciseRow[] => {
  const next = structuredClone(rows);
  if (next[rowIndex].sets[0]) {
    next[rowIndex].sets[0].note = note || null;
  }
  return next;
};

/**
 * Exercise名を変更
 */
export const updateExerciseName = (
  rows: ExerciseRow[],
  rowIndex: number,
  name: string,
): ExerciseRow[] => {
  const next = structuredClone(rows);
  next[rowIndex].name = name;
  // 自由入力なので既存IDに紐づけない（新規扱い）
  next[rowIndex].id = name.trim() ? (next[rowIndex].id ?? null) : null;
  next[rowIndex] = ensureFiveSets(next[rowIndex]);
  return next;
};

// ==================== Form Data Transformers ====================

/**
 * フォームデータを送信用に変換
 * 空のセット（weight_kg と reps が両方とも空）は除外
 */
export const transformFormDataForSubmit = (
  data: WorkoutFormDTO,
  year: number,
  month: number,
  day: number,
) => ({
  ...data,
  performed_date: formatDate(year, month, day),
  started_at: data.started_at || null, // HH:mm形式のまま送る
  ended_at: data.ended_at || null, // HH:mm形式のまま送る
  gym_name: data.gym_name?.trim() || null, // ジム名をバックエンドに送信
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
      .filter((s) => s.weight_kg !== null || s.reps !== null), // 空のセットを除外
  })),
});
