// app/(models)/training_set.ts
export type TrainingSet = {
  id: number;
  trainingRecordId: number;
  trainingExerciseId: number;
  setNumber: number;
  weightKg: number;
  reps: number;
  estimatedMax: number | null;
  note: string | null;
  createdAt: string;
  updatedAt: string;
  deletedAt: string | null;
};
