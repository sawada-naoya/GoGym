import { z } from "zod";

const workoutSetSchema = z.object({
  id: z.number().nullable().optional(),
  set_number: z.number(),
  weight_kg: z.union([z.number(), z.string()]),
  reps: z.union([z.number(), z.string()]),
  note: z.string().nullable(),
});

const workoutExerciseSchema = z.object({
  id: z.number().nullable().optional(),
  name: z.string(),
  workout_part_id: z.number().nullable(),
  sets: z.array(workoutSetSchema),
});

export const workoutFormSchema = z.object({
  id: z.number().nullable().optional(),
  performed_date: z.string(),
  started_at: z.string().nullable(),
  ended_at: z.string().nullable(),
  gym_id: z.number().nullable(),
  gym_name: z.string().nullable(),
  note: z.string().nullable(),
  condition_level: z
    .union([
      z.literal(1),
      z.literal(2),
      z.literal(3),
      z.literal(4),
      z.literal(5),
    ])
    .nullable(),
  exercises: z.array(workoutExerciseSchema),
});

// 型チェック: ZodスキーマがWorkoutFormDTOと一致することを保証
import type { WorkoutFormDTO } from "@/types/workout";
type _Check =
  WorkoutFormDTO extends z.infer<typeof workoutFormSchema>
    ? z.infer<typeof workoutFormSchema> extends WorkoutFormDTO
      ? true
      : false
    : false;
