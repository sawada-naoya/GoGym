// ==================== API Response Types ====================

export type WorkoutRecord = {
  id: number;
  user_id: string;
  performed_date: string; // "YYYY-MM-DD"
  started_at: string | null; // "YYYY-MM-DDTHH:mm:ssZ" など
  ended_at: string | null;
  place: string | null;
  note: string | null;
  condition_level: 1 | 2 | 3 | 4 | 5 | null;
  duration_minutes: number | null;

  WorkoutParts?: WorkoutPart[];
};

export type WorkoutPart = {
  id: number;
  name: string;
  user_id: string | null;

  WorkoutExercises?: WorkoutExercise[];
};

export type WorkoutExercise = {
  id: number;
  name: string;
  workout_part_id: number | null;
  user_id: string | null;

  sets: WorkoutSet[];
};

export type WorkoutSet = {
  id: number;
  workout_record_id: number;
  workout_exercise_id: number;
  set_number: number;
  weight_kg: number;
  reps: number;
  estimated_max: number | null;
  note: string | null;
};

// ==================== Form Types ====================

export type WorkoutPartDTO = {
  id: number;
  key: string;
  translations: Array<{
    locale: string;
    name: string;
  }>;
  exercises: Array<{
    id: number;
    name: string;
    workout_part_id: number | null;
  }>;
};

export type GymDTO = {
  id: number;
  name: string;
};

export type WorkoutFormDTO = {
  id?: number | null;
  performed_date: string; // "YYYY-MM-DD"
  started_at: string | null; // "HH:mm"
  ended_at: string | null; // "HH:mm"
  gym_id: number | null; // deprecated
  gym_name: string | null;
  note: string | null;
  condition_level: 1 | 2 | 3 | 4 | 5 | null;

  exercises: {
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
  }[];
};

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

// バックエンドから返ってくるレスポンス型
export type WorkoutRecordResponseDTO = {
  id?: number | null;
  performed_date: string;
  started_at?: string | null;
  ended_at?: string | null;
  gym_id?: number | null;
  gym_name?: string | null;
  note?: string | null;
  condition_level?: 1 | 2 | 3 | 4 | 5 | null;
  parts: Array<{
    id: number;
    key: string;
    translations: Array<{
      locale: string;
      name: string;
    }>;
    exercises: Array<{
      id?: number | null;
      name: string;
      workout_part_id?: number | null;
      sets: Array<{
        id?: number | null;
        set_number: number;
        weight_kg?: number | null;
        reps?: number | null;
        note?: string | null;
      }>;
    }>;
  }>;
};
