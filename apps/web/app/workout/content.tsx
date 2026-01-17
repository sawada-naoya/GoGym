"use client";
import { useTranslation } from "react-i18next";
import { useBanner } from "@/components/Banner";
import { transformFormDataForSubmit } from "@/features/workout/lib/transforms";
import { useWorkoutForm } from "@/features/workout/hooks/useWorkoutForm";
import { FormProvider, useWatch } from "react-hook-form";
import WorkoutMetadataEditor from "@/features/workout/components/WorkoutMetadataEditor";
import WorkoutExercisesEditor from "@/features/workout/components/WorkoutExercisesEditor";
import { useWorkoutDate } from "@/features/workout/hooks/useWorkoutDate";
import { useEffect, useState } from "react";
import {
  createWorkoutRecord,
  updateWorkoutRecord,
  getWorkoutParts,
} from "@/features/workout/actions";
import type { WorkoutFormDTO, WorkoutPartDTO } from "@/types/workout";
import { useLastRecords } from "@/features/workout/hooks/useLastRecords";

type Props = {
  defaultValues: WorkoutFormDTO;
  availableParts: WorkoutPartDTO[];
  isUpdate: boolean;
};

const WorkoutContent = ({
  defaultValues,
  availableParts: initialParts,
  isUpdate,
}: Props) => {
  const { t } = useTranslation("common");
  const { success, error } = useBanner();

  const { year, month, day } = useWorkoutDate();
  const { fetchLastRecord } = useLastRecords();

  const [availableParts, setAvailableParts] =
    useState<WorkoutPartDTO[]>(initialParts);

  const { form, handleSubmit, isSubmitting } = useWorkoutForm({
    defaultValues,
    onSubmit: async (data) => {
      const body = transformFormDataForSubmit(data, year, month, day);

      try {
        if (data.id) {
          const result = await updateWorkoutRecord(String(data.id), body);
          if (!result.success) {
            return error(t("workout.exercises.errorUpdateFailed"));
          }
          success(t("workout.exercises.successUpdate"));
        } else {
          const result = await createWorkoutRecord(body);
          if (!result.success) {
            return error(t("workout.exercises.errorSaveFailed"));
          }
          success(t("workout.exercises.successSave"));
        }
      } catch {
        error(t("workout.exercises.errorNetworkError"));
      }
    },
  });

  const allExercises = useWatch({ control: form.control, name: "exercises" });

  const [selectedPartId, setSelectedPartId] = useState<number | null>(() => {
    return defaultValues.exercises?.[0]?.workout_part_id ?? null;
  });

  useEffect(() => {
    const newPartId = defaultValues.exercises?.[0]?.workout_part_id ?? null;
    setSelectedPartId(newPartId);
  }, [defaultValues.exercises]);

  const refetchWorkoutParts = async () => {
    try {
      const result = await getWorkoutParts();
      if (result.success && result.data) {
        setAvailableParts(result.data);
      }
    } catch (err) {
      console.error("Failed to refetch workout parts:", err);
    }
  };

  return (
    <FormProvider {...form}>
      <div className="min-h-screen bg-gray-50 py-8">
        <div className="container mx-auto px-4 max-w-7xl">
          <div className="flex items-center justify-between mb-4">
            <div className="text-xl font-semibold">{t("workout.title")}</div>
          </div>

          {/* メタデータエディター（日付・時間・場所・コンディション） */}
          <WorkoutMetadataEditor />

          {/* 種目エディター（部位選択・種目・セット） */}
          <WorkoutExercisesEditor
            allExercises={allExercises ?? defaultValues.exercises}
            onChangeExercises={(next) =>
              form.setValue("exercises", next, { shouldDirty: true })
            }
            workoutParts={availableParts}
            selectedPartId={selectedPartId}
            onPartChange={setSelectedPartId}
            isUpdate={isUpdate}
            onSubmit={handleSubmit}
            onRefetchParts={refetchWorkoutParts}
            dataKey={`${year}-${month}-${day}`}
            onFetchLastRecord={fetchLastRecord}
          />
        </div>
      </div>
      <button onClick={handleSubmit} disabled={isSubmitting}>
        {isUpdate
          ? t("workout.exercises.updateButton")
          : t("workout.exercises.registerButton")}
      </button>
    </FormProvider>
  );
};

export default WorkoutContent;
