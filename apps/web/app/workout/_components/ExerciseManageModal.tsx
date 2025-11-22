"use client";
import React, { useState } from "react";
import { useSession } from "next-auth/react";
import type { WorkoutPartDTO } from "../_lib/types";
import { upsertWorkoutExercises } from "../_lib/api";

type Props = {
  isOpen: boolean;
  onClose: () => void;
  workoutParts: WorkoutPartDTO[];
  onSuccess?: () => void;
};

const ExerciseManageModal: React.FC<Props> = ({ isOpen, onClose, workoutParts, onSuccess }) => {
  const { data: session } = useSession();
  const token = session?.user?.accessToken || "";
  const [selectedPart, setSelectedPart] = useState<number | null>(null);
  const [exerciseNames, setExerciseNames] = useState<string[]>([""]);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // モーダルが開いたときにリセット
  React.useEffect(() => {
    if (isOpen) {
      setExerciseNames([""]);
      setSelectedPart(null);
    }
  }, [isOpen]);

  if (!isOpen) return null;

  const handleAddExerciseInput = () => {
    setExerciseNames([...exerciseNames, ""]);
  };

  const handleUpdateExerciseName = (index: number, value: string) => {
    const updated = [...exerciseNames];
    updated[index] = value;
    setExerciseNames(updated);
  };

  const handleRemoveExercise = (index: number) => {
    if (exerciseNames.length > 1) {
      setExerciseNames(exerciseNames.filter((_, i) => i !== index));
    }
  };

  const handleRegister = async () => {
    if (!selectedPart) {
      alert("部位を選択してください");
      return;
    }

    const validNames = exerciseNames.filter((name) => name.trim() !== "");
    if (validNames.length === 0) {
      alert("種目名を入力してください");
      return;
    }

    setIsSubmitting(true);
    try {
      const exercisesWithPart = validNames.map((name) => ({
        name,
        workout_part_id: selectedPart,
      }));

      const result = await upsertWorkoutExercises(token, exercisesWithPart);

      if (result.ok) {
        alert("種目を登録しました");
        onSuccess?.();
        onClose();
      } else {
        alert(result.error || "登録に失敗しました");
      }
    } catch (error) {
      console.error("種目登録エラー:", error);
      alert("登録中にエラーが発生しました");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
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

        {/* 新規種目追加フォーム */}
        <div className="flex-1 px-6 py-4 overflow-y-auto">
          <div className="space-y-3">
            {exerciseNames.map((name, index) => (
              <div key={index} className="flex gap-2">
                <input type="text" value={name} onChange={(e) => handleUpdateExerciseName(index, e.target.value)} placeholder="種目名" className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" />
                <button type="button" onClick={() => handleRemoveExercise(index)} className="p-2 text-gray-400 hover:text-red-600 transition-colors" disabled={exerciseNames.length === 1}>
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
  );
};

export default ExerciseManageModal;
