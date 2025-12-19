"use client";
import { useEffect, useState } from "react";
import { useForm, FormProvider, useWatch } from "react-hook-form";
import WorkoutMetadataEditor from "./_components/WorkoutMetadataEditor";
import WorkoutExercisesEditor from "./_components/WorkoutExercisesEditor";
import { WorkoutFormDTO, WorkoutPartDTO, ExerciseDTO } from "./_lib/types";
import { transformFormDataForSubmit } from "./_lib/utils";
import { useBanner } from "@/components/Banner";
import {
  getWorkoutParts,
  createWorkoutRecord,
  updateWorkoutRecord,
  getLastExerciseRecord,
} from "@/lib/bff/workout";

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
        const result = await updateWorkoutRecord(data.id, body);
        if (!result.ok) return error(result.error || "更新に失敗しました");
        success("更新しました");
      } else {
        const result = await createWorkoutRecord(body);
        if (!result.ok) return error(result.error || "保存に失敗しました");
        success("保存しました");
      }
    } catch {
      error("通信エラーが発生しました");
    }
  };

  const refetchWorkoutParts = async () => {
    try {
      const parts = await getWorkoutParts();
      setAvailableParts(parts);
    } catch (err) {
      console.error("Failed to refetch workout parts:", err);
    }
  };

  return (
    <FormProvider {...form}>
      <div className="min-h-screen bg-gray-50 py-8">
        <div className="container mx-auto px-4 max-w-7xl">
          <div className="flex items-center justify-between mb-4">
            <div className="text-xl font-semibold">トレーニングノート</div>
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
            onFetchLastRecord={getLastExerciseRecord}
          />
        </div>
      </div>
    </FormProvider>
  );
};

export default WorkoutContent;
