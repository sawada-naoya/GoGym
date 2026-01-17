export const validationMessages = {
  ja: {
    exerciseName: {
      required: "種目名を入力してください",
    },
    sets: {
      min: "最低1セット必要です",
    },
    exercises: {
      min: "最低1つの種目が必要です",
    },
    date: {
      format: "日付の形式が正しくありません",
    },
    time: {
      format: "時刻の形式が正しくありません (HH:mm)",
      invalid: "開始時刻は終了時刻より前にしてください",
    },
    set: {
      bothRequired: "重量と回数は両方入力してください",
    },
  },
  en: {
    exerciseName: {
      required: "Exercise name is required",
    },
    sets: {
      min: "At least 1 set is required",
    },
    exercises: {
      min: "At least 1 exercise is required",
    },
    date: {
      format: "Invalid date format",
    },
    time: {
      format: "Invalid time format (HH:mm)",
      invalid: "Start time must be before end time",
    },
    set: {
      bothRequired: "Both weight and reps are required",
    },
  },
} as const;

/**
 * 現在のロケールに応じたメッセージを取得
 */
export const getValidationMessages = (locale: "ja" | "en" = "ja") => {
  return validationMessages[locale];
};
