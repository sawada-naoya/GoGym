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
  name: string;
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
