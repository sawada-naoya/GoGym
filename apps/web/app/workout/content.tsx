"use client";
import { useEffect, useState } from "react";
import { useForm, FormProvider, useWatch } from "react-hook-form";
import WorkoutMetadataEditor from "./_components/WorkoutMetadataEditor";
import WorkoutExercisesEditor from "./_components/WorkoutExercisesEditor";
import { WorkoutFormDTO, WorkoutPartDTO } from "./_lib/types";
import { transformFormDataForSubmit } from "./_lib/utils";
import { createWorkoutRecord, updateWorkoutRecord, fetchWorkoutRecord, fetchWorkoutParts } from "./_lib/api";
import { useBanner } from "@/components/Banner";
import { formatDate } from "@/lib/time";

type Props = {
  Year: number;
  Month: number;
  Day: number;
  defaultValues: WorkoutFormDTO;
  availableParts: WorkoutPartDTO[];
  isUpdate: boolean;
  token: string;
};

const WorkoutContent = ({ Year, Month, Day, defaultValues, availableParts: initialParts, isUpdate, token }: Props) => {
  const { success, error } = useBanner();
  const [selectedDay, setSelectedDay] = useState(Day);
  const [selectedYear, setSelectedYear] = useState(Year);
  const [selectedMonth, setSelectedMonth] = useState(Month);
  const [isLoadingPart, setIsLoadingPart] = useState(false);
  const [availableParts, setAvailableParts] = useState<WorkoutPartDTO[]>(initialParts);

  const form = useForm<WorkoutFormDTO>({
    defaultValues,
    mode: "onBlur",
  });

  const rows = useWatch({ control: form.control, name: "exercises" });
  const selectedPart = useWatch({ control: form.control, name: "workout_part" });

  // Props の Year, Month, Day が変わったら state を同期
  useEffect(() => {
    setSelectedYear(Year);
    setSelectedMonth(Month);
    setSelectedDay(Day);
  }, [Year, Month, Day]);

  useEffect(() => {
    form.reset(defaultValues);
  }, [defaultValues, form]);

  // 部位変更時にレコードを再取得
  useEffect(() => {
    if (!selectedPart?.id) return;

    const fetchPartRecords = async () => {
      setIsLoadingPart(true);
      try {
        const date = formatDate(selectedYear, selectedMonth, selectedDay);
        const data = await fetchWorkoutRecord(token, date, selectedPart.id);
        // 部位とexercisesのみ更新（メタデータは維持）
        form.setValue("exercises", data.exercises, { shouldDirty: false });
      } catch (err) {
        console.error("Failed to fetch part records:", err);
      } finally {
        setIsLoadingPart(false);
      }
    };

    fetchPartRecords();
  }, [selectedPart?.id, selectedYear, selectedMonth, selectedDay, token, form]);

  const handleSubmit = async (data: WorkoutFormDTO) => {
    const body = transformFormDataForSubmit(data, selectedYear, selectedMonth, selectedDay);

    try {
      if (isUpdate && data.id) {
        const result = await updateWorkoutRecord(token, data.id, body);
        if (!result.ok) return error(result.error || "更新に失敗しました");
        success("更新しました");
      } else {
        const result = await createWorkoutRecord(token, body);
        if (!result.ok) return error(result.error || "保存に失敗しました");
        success("保存しました");
      }
    } catch {
      error("通信エラーが発生しました");
    }
  };

  const refetchWorkoutParts = async () => {
    try {
      const parts = await fetchWorkoutParts(token);
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
          <WorkoutMetadataEditor form={form} selectedYear={selectedYear} selectedMonth={selectedMonth} selectedDay={selectedDay} onYearChange={setSelectedYear} onMonthChange={setSelectedMonth} onDayChange={setSelectedDay} />

          {/* 種目エディター（部位選択・種目・セット） */}
          <WorkoutExercisesEditor
            workoutExercises={rows ?? defaultValues.exercises}
            onChangeExercises={(next) => form.setValue("exercises", next, { shouldDirty: true })}
            workoutParts={availableParts}
            selectedPart={form.watch("workout_part")}
            onPartChange={(part) => form.setValue("workout_part", part, { shouldDirty: true })}
            isUpdate={isUpdate}
            onSubmit={form.handleSubmit(handleSubmit)}
            onRefetchParts={refetchWorkoutParts}
          />
        </div>
      </div>
    </FormProvider>
  );
};

export default WorkoutContent;
