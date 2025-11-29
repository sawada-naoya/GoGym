"use client";
import { useState, useEffect } from "react";
import { useSession } from "next-auth/react";
import type { WorkoutPartDTO } from "../_lib/types";
import { upsertWorkoutExercises, deleteWorkoutExercise } from "../_lib/api";
import { useBanner } from "@/components/Banner";

type Props = {
  isOpen: boolean;
  onClose: () => void;
  workoutParts: WorkoutPartDTO[];
  onSuccess?: () => void;
};

type ExerciseFormItem = {
  id?: number;
  name: string;
};

const ExerciseManageModal: React.FC<Props> = ({ isOpen, onClose, workoutParts, onSuccess }) => {
  const { data: session } = useSession();
  const token = session?.user?.accessToken || "";
  const { success, error } = useBanner();
  const [selectedPart, setSelectedPart] = useState<number | null>(null);
  const [exercises, setExercises] = useState<ExerciseFormItem[]>([{ name: "" }]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState<{ index: number; exercise: ExerciseFormItem } | null>(null);

  // 部位選択時に既存種目をフォームにセット
  useEffect(() => {
    if (!selectedPart) {
      setExercises([{ name: "" }]);
      return;
    }

    const part = workoutParts.find((p) => p.id === selectedPart);
    const existing = part?.exercises || [];

    if (existing.length > 0) {
      setExercises(existing.map((ex) => ({ id: ex.id, name: ex.name })));
    } else {
      setExercises([{ name: "" }]);
    }
  }, [selectedPart, workoutParts]);

  // モーダルが開いたときにリセット
  useEffect(() => {
    if (isOpen) {
      setExercises([{ name: "" }]);
      setSelectedPart(null);
    }
  }, [isOpen]);

  if (!isOpen) return null;

  const handleAddExerciseInput = () => {
    setExercises([...exercises, { name: "" }]);
  };

  const handleUpdateExerciseName = (index: number, value: string) => {
    const updated = [...exercises];
    updated[index].name = value;
    setExercises(updated);
  };

  const handleRemoveExercise = (index: number) => {
    const exercise = exercises[index];

    // 新規（IDなし）の場合は確認なしで削除
    if (!exercise.id) {
      if (exercises.length > 1) {
        setExercises(exercises.filter((_, i) => i !== index));
      }
      return;
    }

    // 既存（IDあり）の場合は確認モーダル表示
    setDeleteTarget({ index, exercise });
  };

  const handleConfirmDelete = async () => {
    if (!deleteTarget || !deleteTarget.exercise.id) return;

    setIsSubmitting(true);
    try {
      const result = await deleteWorkoutExercise(token, deleteTarget.exercise.id);

      if (result.ok) {
        // 削除成功：フロント側も削除
        if (exercises.length > 1) {
          setExercises(exercises.filter((_, i) => i !== deleteTarget.index));
        }
        success("種目を削除しました");
        onSuccess?.(); // 親側で再取得させる
      } else {
        error(result.error || "削除に失敗しました");
      }
    } catch (err) {
      console.error("種目削除エラー:", err);
      error("削除中にエラーが発生しました");
    } finally {
      setIsSubmitting(false);
      setDeleteTarget(null);
    }
  };

  const handleCancelDelete = () => {
    setDeleteTarget(null);
  };

  const handleRegister = async () => {
    if (!selectedPart) {
      error("部位を選択してください");
      return;
    }

    const validExercises = exercises.filter((ex) => ex.name.trim() !== "");
    if (validExercises.length === 0) {
      error("種目名を入力してください");
      return;
    }

    setIsSubmitting(true);
    try {
      const exercisesWithPart = validExercises.map((ex) => ({
        id: ex.id,
        name: ex.name,
        workout_part_id: selectedPart,
      }));

      const result = await upsertWorkoutExercises(token, exercisesWithPart);

      if (result.ok) {
        success("種目を登録しました");
        onSuccess?.();
        // モーダルは閉じない
      } else {
        error(result.error || "登録に失敗しました");
      }
    } catch (err) {
      console.error("種目登録エラー:", err);
      error("登録中にエラーが発生しました");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <>
      {/* 削除確認モーダル */}
      {deleteTarget && (
        <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-[60] p-4">
          <div className="bg-white rounded-lg shadow-xl w-[400px] p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-2">種目を削除しますか？</h3>
            <p className="text-gray-600 mb-6">「{deleteTarget.exercise.name}」を削除します。この操作は取り消せません。</p>
            <div className="flex gap-3 justify-center">
              <button onClick={handleCancelDelete} className="px-4 py-2 rounded-md border border-gray-300 text-gray-700 hover:bg-gray-50 transition-colors">
                キャンセル
              </button>
              <button onClick={handleConfirmDelete} className="px-4 py-2 rounded-md bg-red-600 text-white hover:bg-red-700 transition-colors">
                削除する
              </button>
            </div>
          </div>
        </div>
      )}

      {/* メインモーダル */}
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
        <div className="bg-white rounded-lg shadow-xl w-[500px] h-[500px] flex flex-col">
          {/* ヘッダー */}
          <div className="px-6 py-4 border-b border-gray-200">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <select value={selectedPart?.toString() ?? ""} onChange={(e) => setSelectedPart(e.target.value ? Number(e.target.value) : null)} className="w-32 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500 bg-white text-base">
                  <option value="">部位</option>
                  {workoutParts.map((part) => (
                    <option key={part.id} value={part.id.toString()}>
                      {part.name}
                    </option>
                  ))}
                </select>
                <h2 className="text-xl font-semibold text-gray-900">の種目を追加</h2>
              </div>
              <button onClick={onClose} className="text-gray-400 hover:text-gray-600 transition-colors">
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>

          {/* 種目フォーム（既存 + 新規） */}
          <div className="flex-1 px-6 py-4 overflow-y-auto">
            <div className="space-y-3">
              {exercises.map((exercise, index) => (
                <div key={index} className="flex gap-2">
                  <input type="text" value={exercise.name} onChange={(e) => handleUpdateExerciseName(index, e.target.value)} placeholder="種目名" className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" />
                  <button type="button" onClick={() => handleRemoveExercise(index)} className="p-2 text-gray-400 hover:text-red-600 transition-colors" disabled={exercises.length === 1}>
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              ))}
              <div className="flex justify-center">
                <button type="button" onClick={handleAddExerciseInput} className="flex items-center justify-center w-10 h-10 text-booking-600 hover:text-booking-700 hover:bg-booking-50 rounded-full transition-colors">
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                  </svg>
                </button>
              </div>
            </div>
          </div>

          {/* フッター */}
          <div className="px-6 py-4 flex justify-center">
            <button onClick={handleRegister} disabled={isSubmitting} className="px-8 py-2 rounded-md bg-booking-600 text-white hover:bg-booking-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
              {isSubmitting ? "登録中..." : "登録"}
            </button>
          </div>
        </div>
      </div>
    </>
  );
};

export default ExerciseManageModal;
