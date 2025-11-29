// app/(models)/training_record.ts
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
