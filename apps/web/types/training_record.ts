// app/(models)/training_record.ts
export type TrainingRecord = {
  id: number;
  userId: string;
  performedDate: string; // "YYYY-MM-DD"
  startedAt: string | null; // "YYYY-MM-DDTHH:mm:ssZ" など
  endedAt: string | null;
  place: string | null;
  note: string | null;
  conditionLevel: 1 | 2 | 3 | 4 | 5 | null;
  durationMinutes: number | null;
  createdAt: string;
  updatedAt: string;
  deletedAt: string | null;
};
