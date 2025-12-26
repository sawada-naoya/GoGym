"use client";
import { useEffect, useState } from "react";
import { useForm, FormProvider, useWatch } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { useBanner } from "@/components/Banner";
import { WorkoutFormDTO, WorkoutPartDTO, ExerciseDTO } from "@/types/workout";
import { transformFormDataForSubmit } from "@/features/workout/lib/utils";
import WorkoutMetadataEditor from "@/features/workout/components/WorkoutMetadataEditor";
import WorkoutExercisesEditor from "@/features/workout/components/WorkoutExercisesEditor";

type Props = {
  Year: number;
  Month: number;
  Day: number;
  defaultValues: WorkoutFormDTO;
  availableParts: WorkoutPartDTO[];
  isUpdate: boolean;
};

const WorkoutContent = ({
  Year,
  Month,
  Day,
  defaultValues,
  availableParts: initialParts,
  isUpdate,
}: Props) => {
  const { t } = useTranslation("common");
  const { success, error } = useBanner();
  const [selectedDay, setSelectedDay] = useState(Day);
  const [selectedYear, setSelectedYear] = useState(Year);
  const [selectedMonth, setSelectedMonth] = useState(Month);
  const [availableParts, setAvailableParts] =
    useState<WorkoutPartDTO[]>(initialParts);

  // 既存データがある場合、最初のexerciseの部位を初期選択
  const initialPartId = defaultValues.exercises?.[0]?.workout_part_id ?? null;
  const [selectedPartId, setSelectedPartId] = useState<number | null>(
    initialPartId,
  );

  const form = useForm<WorkoutFormDTO>({
    defaultValues,
    mode: "onBlur",
  });

  const allExercises = useWatch({ control: form.control, name: "exercises" });

  // Props の Year, Month, Day が変わったら state を同期
  useEffect(() => {
    setSelectedYear(Year);
    setSelectedMonth(Month);
    setSelectedDay(Day);
  }, [Year, Month, Day]);

  useEffect(() => {
    form.reset(defaultValues);
    // defaultValuesが変わったら部位も再設定
    const newPartId = defaultValues.exercises?.[0]?.workout_part_id ?? null;
    setSelectedPartId(newPartId);
  }, [defaultValues]);

  const handleSubmit = async (data: WorkoutFormDTO) => {
    const body = transformFormDataForSubmit(
      data,
      selectedYear,
      selectedMonth,
      selectedDay,
    );

    try {
      if (isUpdate && data.id) {
        const res = await fetch(`/api/workouts/records/${data.id}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(body),
          cache: "no-store",
        });
        if (!res.ok) return error(t("workout.exercises.errorUpdateFailed"));
        success(t("workout.exercises.successUpdate"));
      } else {
        const res = await fetch("/api/workouts/records", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(body),
          cache: "no-store",
        });
        if (!res.ok) return error(t("workout.exercises.errorSaveFailed"));
        success(t("workout.exercises.successSave"));
      }
    } catch {
      error(t("workout.exercises.errorNetworkError"));
    }
  };

  const refetchWorkoutParts = async () => {
    try {
      const res = await fetch("/api/workouts/parts", { cache: "no-store" });
      if (res.ok) {
        const parts = await res.json();
        setAvailableParts(parts);
      }
    } catch (err) {
      console.error("Failed to refetch workout parts:", err);
    }
  };

  const fetchLastExerciseRecord = async (
    exerciseID: number,
  ): Promise<ExerciseDTO | null> => {
    if (!exerciseID) {
      return null;
    }

    try {
      const res = await fetch(
        `/api/workouts/exercises?id=${exerciseID}&action=last`,
        { cache: "no-store" },
      );

      if (!res.ok || res.status === 204) {
        return null;
      }

      const data = await res.json();
      return data || null;
    } catch {
      return null;
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
          <WorkoutMetadataEditor
            form={form}
            selectedYear={selectedYear}
            selectedMonth={selectedMonth}
            selectedDay={selectedDay}
            onYearChange={setSelectedYear}
            onMonthChange={setSelectedMonth}
            onDayChange={setSelectedDay}
          />

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
            onSubmit={form.handleSubmit(handleSubmit)}
            onRefetchParts={refetchWorkoutParts}
            dataKey={`${selectedYear}-${selectedMonth}-${selectedDay}`}
            onFetchLastRecord={fetchLastExerciseRecord}
          />
        </div>
      </div>
    </FormProvider>
  );
};

export default WorkoutContent;
