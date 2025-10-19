export type WorkoutFormDTO = {
  id?: number | null; // 既存ならrecord id
  performedDate: string; // "YYYY-MM-DD"
  startedAt: string | null; // "HH:mm"（フォーム用）
  endedAt: string | null; // "HH:mm"
  place: string;
  note: string | null;
  conditionLevel: 1 | 2 | 3 | 4 | 5 | null;

  workoutPart: {
    id: number | null;
    name: string | null;
    source: "preset" | "custom" | null;
  };

  exercises: {
    id?: number | null; // workout_exercises.id（既存なら）
    name: string; // 種目名
    workoutPartId: number | null;
    isDefault: 0 | 1;
    sets: {
      id?: number | null; // workout_sets.id（既存なら）
      setNumber: number;
      weightKg: number | string; // 入力中は空文字も許容
      reps: number | string;
      note: string | null; // 入力中は空文字も許容
    }[];
  }[];
};
