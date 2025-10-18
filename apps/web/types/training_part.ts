export type TrainingPart = {
  id: number;
  name: string;
  isDefault: 0 | 1;
  userId: string | null;
  createdAt: string;
  updatedAt: string;
  deletedAt: string | null;
};
