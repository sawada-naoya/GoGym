export type WorkoutPartDTO = {
  id: number;
  name: string;
  isDefault: boolean;
};

export type WorkoutFormDTO = {
  id?: number | null; // 既存ならrecord id
  performed_date: string; // "YYYY-MM-DD"
  started_at: string | null; // "HH:mm"（フォーム用）
  ended_at: string | null; // "HH:mm"
  place: string;
  note: string | null;
  condition_level: 1 | 2 | 3 | 4 | 5 | null;

  workout_part: {
    id: number | null;
    name: string | null;
    source: "preset" | "custom" | null;
  };

  exercises: {
    id?: number | null; // workout_exercises.id（既存なら）
    name: string; // 種目名
    workout_part_id: number | null;
    is_default: 0 | 1;
    sets: {
      id?: number | null; // workout_sets.id（既存なら）
      set_number: number;
      weight_kg: number | string; // 入力中は空文字も許容
      reps: number | string;
      note: string | null; // 入力中は空文字も許容
    }[];
  }[];
};
