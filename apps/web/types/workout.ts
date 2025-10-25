// app/(models)/training_record.ts
export type WorkoutRecord = {
  id: number;
  user_id: string;
  performed_date: string; // "YYYY-MM-DD"
  startedAt: string | null; // "YYYY-MM-DDTHH:mm:ssZ" など
  endedAt: string | null;
  place: string | null;
  note: string | null;
  condition_level: 1 | 2 | 3 | 4 | 5 | null;
  duration_minutes: number | null;
};

export type WorkoutPart = {
  id: number;
  name: string;
  is_default: 0 | 1;
  user_id: string | null;
};

export type TrainingExercise = {
  id: number;
  name: string;
  workout_part_id: number | null;
  is_default: 0 | 1;
  user_id: string | null;
};

export type WorkoutSet = {
  id: number;
  workout_record_id: number;
  training_exerciseId: number;
  set_number: number;
  weight_kg: number;
  reps: number;
  estimated_max: number | null;
  note: string | null;
};
