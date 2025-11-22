"use client";
import { useEffect, useState } from "react";
import { useForm, FormProvider, useWatch } from "react-hook-form";
import WorkoutMetadataEditor from "./_components/WorkoutMetadataEditor";
import WorkoutExercisesEditor from "./_components/WorkoutExercisesEditor";
import { WorkoutFormDTO, WorkoutPartDTO } from "./_lib/types";
import { transformFormDataForSubmit } from "./_lib/utils";
import { createWorkoutRecord, updateWorkoutRecord } from "./_lib/api";
import { useBanner } from "@/components/Banner";

type Props = {
  Year: number;
  Month: number;
  Day: number;
  defaultValues: WorkoutFormDTO;
  availableParts: WorkoutPartDTO[];
  isUpdate: boolean;
  token: string;
};

const WorkoutContent = ({ Year, Month, Day, defaultValues, availableParts, isUpdate, token }: Props) => {
  const { success, error } = useBanner();
  const [selectedDay, setSelectedDay] = useState(Day);
  const [selectedYear, setSelectedYear] = useState(Year);
  const [selectedMonth, setSelectedMonth] = useState(Month);

  const form = useForm<WorkoutFormDTO>({
    defaultValues,
    mode: "onBlur",
  });

  const rows = useWatch({ control: form.control, name: "exercises" });

  useEffect(() => {
    form.reset(defaultValues);
  }, [defaultValues, form]);

  const handleSubmit = async (data: WorkoutFormDTO) => {
    const body = transformFormDataForSubmit(data, Year, Month, selectedDay);

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
          />
        </div>
      </div>
    </FormProvider>
  );
};

export default WorkoutContent;
