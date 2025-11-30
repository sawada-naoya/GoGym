export type WorkoutPartDTO = {
  id: number;
  name: string;
  exercises: Array<{
    id: number;
    name: string;
    workout_part_id: number | null;
  }>;
};

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
  exercises: Array.from({ length: 1 }, () => ({
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
  })),
});

export type ExerciseDTO = {
  id?: number | null;
  name: string;
  workout_part_id: number | null;
  sets: {
    id?: number | null;
    set_number: number;
    weight_kg: number | string;
    reps: number | string;
    note: string | null;
  }[];
};

export type WorkoutFormDTO = {
  id?: number | null; // 既存ならrecord id
  performed_date: string; // "YYYY-MM-DD"
  started_at: string | null; // "HH:mm"（フォーム用）
  ended_at: string | null; // "HH:mm"
  place: string;
  note: string | null;
  condition_level: 1 | 2 | 3 | 4 | 5 | null;

  exercises: {
    id?: number | null; // workout_exercises.id（既存なら）
    name: string; // 種目名
    workout_part_id: number | null;
    sets: {
      id?: number | null; // workout_sets.id（既存なら）
      set_number: number;
      weight_kg: number | string; // 入力中は空文字も許容
      reps: number | string;
      note: string | null; // 入力中は空文字も許容
    }[];
  }[];
};
