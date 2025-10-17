// 部位の出所
export type PartSource = "preset" | "custom" | null;

// セット入力（フォーム用）
export type TrainingSetInput = {
  setNumber: number;
  weightKg: number | ""; // 未入力は ""（送信前に0へ正規化OK）
  reps: number | "";
  note?: string | null;
};

// 種目＋セット（名前も一緒に持つ）
export type TrainingExerciseWithSetsInput = {
  // 既存を選んだ時は id、カスタムなら name を使用（両方あってもOKだがサーバで id 優先）
  trainingExerciseId?: number;
  trainingExerciseName: string;
  sets: TrainingSetInput[];
};

// 部位（既存を選べば id、カスタムなら name と source="custom"）
export type TrainingPartInput = {
  trainingPartId?: number;
  name: string | null; // "胸" 等。選ばない（null）も許容
  source: PartSource; // "preset" | "custom" | null
};

// 作成/更新のDTO（フォーム送信に使う）
export type TrainingRecordAPI = {
  id?: number; // 更新時のみ
  performedDate: string; // "YYYY-MM-DD"
  startedAt?: string | null; // ISO もしくは "HH:mm"（送信前に正規化）
  endedAt?: string | null;
  place?: string | null;
  note?: string | null;
  conditionLevel?: 1 | 2 | 3 | 4 | 5 | null;

  trainingPart?: TrainingPartInput | null;

  exercises: TrainingExerciseWithSetsInput[];
};

// -----------------------------
// 表示用 (GET /api/training-records/:id)
// -----------------------------

export type TrainingExercise = {
  id: number; // training_exercises.id
  name: string; // "ベンチプレス" など
  trainingPartId: number | null; // 紐づく部位ID
  isDefault: 0 | 1; // プリセット or ユーザー作成
  sets: {
    id: number;
    setNumber: number;
    weightKg: number;
    reps: number;
    estimatedMax: number | null; // 推定1RM（計算済み）
    note: string | null; // メモ
  }[];
};

export type TrainingRecordDisplay = {
  id: number;
  performedDate: string; // "YYYY-MM-DD"
  startedAt: string | null; // ISO
  endedAt: string | null; // ISO
  place: string | null;
  note: string | null;
  conditionLevel: 1 | 2 | 3 | 4 | 5 | null;
  durationMinutes: number | null;

  // 部位情報（null許容）
  trainingPart: {
    id: number | null;
    name: string | null; // "胸"など
    source: "preset" | "custom" | null;
  } | null;

  // 種目＋セット（一覧）
  exercises: TrainingExercise[];
};
