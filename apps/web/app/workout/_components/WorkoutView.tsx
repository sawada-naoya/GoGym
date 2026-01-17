"use client";

import { useTranslation } from "react-i18next";
import { useFormContext } from "react-hook-form";
import { useWorkoutContext } from "../_providers/WorkoutProvider";
import WorkoutMetadataEditor from "@/features/workout/components/WorkoutMetadataEditor";
import WorkoutExercisesEditor from "@/features/workout/components/WorkoutExercisesEditor";
import type { WorkoutFormDTO } from "@/types/workout";

/**
 * WorkoutView
 *
 * 責務:
 * - 純粋なUI表示コンポーネント
 * - Contextから値を取得（propsは0個）
 * - WorkoutMetadataEditor + WorkoutExercisesEditor を配置
 * - ロジックは一切持たない
 */
export const WorkoutView = () => {
  const { t } = useTranslation("common");

  // Contextから値を取得
  const {
    year,
    month,
    day,
    handleSubmit,
    isSubmitting,
    isUpdate,
    allExercises,
    availableParts,
    selectedPartId,
    setSelectedPartId,
    refetchWorkoutParts,
    fetchLastRecord,
  } = useWorkoutContext();

  // FormProviderから form を取得
  const form = useFormContext<WorkoutFormDTO>();

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="container mx-auto px-4 max-w-7xl">
        {/* ヘッダー */}
        <div className="flex items-center justify-between mb-4">
          <div className="text-xl font-semibold">{t("workout.title")}</div>
        </div>

        {/* メタデータエディター（日付・時間・場所・コンディション） */}
        <WorkoutMetadataEditor />

        {/* 種目エディター（部位選択・種目・セット） */}
        <WorkoutExercisesEditor />
      </div>

      {/* Submit button（画面下部固定） */}
      <div className="md:hidden fixed bottom-0 left-0 right-0 p-3 bg-white border-t border-gray-200 shadow-lg z-10">
        <button
          onClick={handleSubmit}
          disabled={isSubmitting}
          className="w-full py-2.5 text-sm rounded-xl bg-booking-600 text-white hover:bg-booking-700 active:bg-booking-800 transition-colors font-bold shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isUpdate
            ? t("workout.exercises.updateButtonMobile")
            : t("workout.exercises.registerButtonMobile")}
        </button>
      </div>
    </div>
  );
};
