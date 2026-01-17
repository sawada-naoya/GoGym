import type { WorkoutFormDTO, WorkoutPartDTO } from "@/types/workout";
import { parseDateString } from "@/utils/time";

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
  const targetLocale =
    locale || (typeof navigator !== "undefined" ? navigator.language : "ja");
  const languageCode = targetLocale.split("-")[0];

  const exactMatch = part.translations.find((t) => t.locale === targetLocale);
  if (exactMatch) return exactMatch.name;

  const languageMatch = part.translations.find(
    (t) => t.locale.split("-")[0] === languageCode,
  );
  if (languageMatch) return languageMatch.name;

  const jaMatch = part.translations.find((t) => t.locale === "ja");
  if (jaMatch) return jaMatch.name;

  return part.translations[0]?.name || part.key;
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
