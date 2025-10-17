export type TrainingExercise = {
  id: number;
  name: string;
  trainingPartId: number | null;
  isDefault: 0 | 1;
  userId: string | null;
  createdAt: string;
  updatedAt: string;
  deletedAt: string | null;
};
