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
    exercises:
      exercises.length > 0
        ? exercises
        : Array.from({ length: 1 }, () => ({
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
