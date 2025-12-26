"use client";
import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { WorkoutPartDTO } from "@/types/workout";
import { useBanner } from "@/components/Banner";
import { getLocalizedPartName } from "@/features/workout/lib/utils";
// Removed lib/bff/workout import - now using direct fetch to Route Handler (BFF)

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

const ExerciseManageModal: React.FC<Props> = ({
  isOpen,
  onClose,
  workoutParts,
  onSuccess,
}) => {
  const { t } = useTranslation("common");
  const { success, error } = useBanner();
  const [selectedPart, setSelectedPart] = useState<number | null>(null);
  const [exercises, setExercises] = useState<ExerciseFormItem[]>([
    { name: "" },
  ]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState<{
    index: number;
    exercise: ExerciseFormItem;
  } | null>(null);

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
      const res = await fetch(
        `/api/workouts/exercises?id=${deleteTarget.exercise.id}`,
        {
          method: "DELETE",
          cache: "no-store",
        },
      );

      if (res.ok) {
        // 削除成功：フロント側も削除
        if (exercises.length > 1) {
          setExercises(exercises.filter((_, i) => i !== deleteTarget.index));
        }
        success(t("workout.exerciseModal.successDelete"));
        onSuccess?.(); // 親側で再取得させる
      } else {
        error(t("workout.exerciseModal.errorDeleteFailed"));
      }
    } catch (err) {
      error(t("workout.exerciseModal.errorDeleting"));
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
      error(t("workout.exerciseModal.errorSelectPart"));
      return;
    }

    const validExercises = exercises.filter((ex) => ex.name.trim() !== "");
    if (validExercises.length === 0) {
      error(t("workout.exerciseModal.errorEnterExerciseName"));
      return;
    }

    setIsSubmitting(true);
    try {
      const exercisesWithPart = validExercises.map((ex) => ({
        id: ex.id,
        name: ex.name,
        workout_part_id: selectedPart,
      }));

      const res = await fetch("/api/workouts/exercises", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ exercises: exercisesWithPart }),
        cache: "no-store",
      });

      if (res.ok) {
        success(t("workout.exerciseModal.successRegister"));
        onSuccess?.();
        // モーダルは閉じない
      } else {
        error(t("workout.exerciseModal.errorRegisterFailed"));
      }
    } catch (err) {
      error(t("workout.exerciseModal.errorRegistering"));
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <>
      {/* 削除確認モーダル */}
      {deleteTarget && (
        <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-[60] p-3 md:p-4">
          <div className="bg-white rounded-lg shadow-xl w-full max-w-[340px] md:max-w-[400px] p-4 md:p-6">
            <h3 className="text-base md:text-lg font-semibold text-gray-900 mb-1.5 md:mb-2">
              {t("workout.exerciseModal.deleteConfirmTitle")}
            </h3>
            <p className="text-sm md:text-base text-gray-600 mb-4 md:mb-6">
              「{deleteTarget.exercise.name}」
              {t("workout.exerciseModal.deleteConfirmMessage")}
            </p>
            <div className="flex gap-2 md:gap-3 justify-center">
              <button
                onClick={handleCancelDelete}
                className="px-3 md:px-4 py-1.5 md:py-2 text-sm md:text-base rounded-md border border-gray-300 text-gray-700 hover:bg-gray-50 transition-colors"
              >
                {t("workout.exerciseModal.cancelButton")}
              </button>
              <button
                onClick={handleConfirmDelete}
                className="px-3 md:px-4 py-1.5 md:py-2 text-sm md:text-base rounded-md bg-red-600 text-white hover:bg-red-700 transition-colors"
              >
                {t("workout.exerciseModal.deleteButton")}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* メインモーダル */}
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-3 md:p-4">
        <div className="bg-white rounded-lg shadow-xl w-full max-w-[380px] md:max-w-[500px] h-[420px] md:h-[500px] flex flex-col">
          {/* ヘッダー */}
          <div className="px-3 md:px-6 py-3 md:py-4 border-b border-gray-200">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-1.5 md:gap-2 flex-1 min-w-0">
                <select
                  value={selectedPart?.toString() ?? ""}
                  onChange={(e) =>
                    setSelectedPart(
                      e.target.value ? Number(e.target.value) : null,
                    )
                  }
                  className="w-24 md:w-32 px-2 md:px-3 py-1.5 md:py-2 text-sm md:text-base border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500 bg-white"
                >
                  <option value="">
                    {t("workout.exerciseModal.partLabel")}
                  </option>
                  {workoutParts.map((part) => (
                    <option key={part.id} value={part.id.toString()}>
                      {getLocalizedPartName(part)}
                    </option>
                  ))}
                </select>
                <h2 className="text-sm md:text-xl font-semibold text-gray-900 truncate">
                  {t("workout.exerciseModal.title")}
                </h2>
              </div>
              <button
                onClick={onClose}
                className="text-gray-400 hover:text-gray-600 transition-colors flex-shrink-0 ml-2"
              >
                <svg
                  className="w-5 h-5 md:w-6 md:h-6"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>
          </div>

          {/* 種目フォーム（既存 + 新規） */}
          <div className="flex-1 px-3 md:px-6 py-3 md:py-4 overflow-y-auto">
            <div className="space-y-2 md:space-y-3">
              {exercises.map((exercise, index) => (
                <div key={index} className="flex gap-1.5 md:gap-2">
                  <input
                    type="text"
                    value={exercise.name}
                    onChange={(e) =>
                      handleUpdateExerciseName(index, e.target.value)
                    }
                    placeholder={t("workout.exerciseModal.exercisePlaceholder")}
                    className="flex-1 px-2 md:px-3 py-1.5 md:py-2 text-sm md:text-base border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500"
                  />
                  <button
                    type="button"
                    onClick={() => handleRemoveExercise(index)}
                    className="p-1.5 md:p-2 text-gray-400 hover:text-red-600 transition-colors"
                    disabled={exercises.length === 1}
                  >
                    <svg
                      className="w-4 h-4 md:w-5 md:h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                      />
                    </svg>
                  </button>
                </div>
              ))}
              <div className="flex justify-center pt-1">
                <button
                  type="button"
                  onClick={handleAddExerciseInput}
                  className="flex items-center justify-center w-8 h-8 md:w-10 md:h-10 text-booking-600 hover:text-booking-700 hover:bg-booking-50 rounded-full transition-colors"
                >
                  <svg
                    className="w-4 h-4 md:w-5 md:h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 4v16m8-8H4"
                    />
                  </svg>
                </button>
              </div>
            </div>
          </div>

          {/* フッター */}
          <div className="px-3 md:px-6 py-3 md:py-4 flex justify-center border-t border-gray-100">
            <button
              onClick={handleRegister}
              disabled={isSubmitting}
              className="px-6 md:px-8 py-1.5 md:py-2 text-sm md:text-base rounded-md bg-booking-600 text-white hover:bg-booking-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed font-medium"
            >
              {isSubmitting
                ? t("workout.exerciseModal.registering")
                : t("workout.exerciseModal.registerButton")}
            </button>
          </div>
        </div>
      </div>
    </>
  );
};

export default ExerciseManageModal;
