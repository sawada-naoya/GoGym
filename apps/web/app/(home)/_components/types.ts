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
// 送受信用を共通化するための型定義
// -----------------------------

// app/_components/types.ts（例）
export type TrainingFormDTO = {
  id?: number | null; // 既存ならrecord id
  performedDate: string; // "YYYY-MM-DD"
  startedAt: string | null; // "HH:mm"（フォーム用）
  endedAt: string | null; // "HH:mm"
  place: string;
  note: string | null;
  conditionLevel: 1 | 2 | 3 | 4 | 5 | null;

  trainingPart: {
    id: number | null;
    name: string | null;
    source: "preset" | "custom" | null;
  };

  exercises: {
    id?: number | null; // training_exercises.id（既存なら）
    name: string; // 種目名
    trainingPartId: number | null;
    isDefault: 0 | 1;
    sets: {
      id?: number | null; // training_sets.id（既存なら）
      setNumber: number;
      weightKg: number | string; // 入力中は空文字も許容
      reps: number | string;
      note: string | null;
    }[];
  }[];
};
