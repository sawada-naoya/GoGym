"use client";
import { useEffect, useState } from "react";
import { useForm, FormProvider, useWatch } from "react-hook-form";
import MonthlyStrip from "./MonthlyStrip";
import TrainingSessionMetaEditor from "./TrainingSessionMetaEditor";
import TrainingExercisesEditor from "./TrainingExercisesEditor";
import { TrainingFormDTO } from "./types";
import { POST, PUT } from "@/lib/api";
import { useBanner } from "@/components/Banner";

type Props = {
  Year: number;
  Month: number;
  Day: number;
  defaultValues: TrainingFormDTO;
  isUpdate: boolean;
};

const pad = (n: number) => (n < 10 ? `0${n}` : `${n}`);
const toISO = (y: number, m: number, d: number, hm?: string | null) => {
  if (!hm) return null;
  const [hh, mm] = hm.split(":").map(Number);
  return new Date(y, m - 1, d, hh, mm, 0).toISOString();
};

const TrainingRecordEditor = ({ Year, Month, Day, defaultValues, isUpdate }: Props) => {
  const { success, error } = useBanner();
  const [selectedDay, setSelectedDay] = useState(Day);

  const form = useForm<TrainingFormDTO>({
    defaultValues,
    mode: "onBlur",
  });

  const rows = useWatch({ control: form.control, name: "exercises" });

  useEffect(() => {
    form.reset(defaultValues);
  }, [defaultValues, form]);

  const performedDate = `${Year}-${pad(Month)}-${pad(selectedDay)}`;

  const handleSubmit = async (data: TrainingFormDTO) => {
    const body = {
      ...data,
      performedDate,
      startedAt: toISO(Year, Month, selectedDay, data.startedAt),
      endedAt: toISO(Year, Month, selectedDay, data.endedAt),
    };

    try {
      if (isUpdate && data.id) {
        const res = await PUT(`/api/training-records/${data.id}`, { body });
        if (!res.ok) return error("更新に失敗しました");
        success("更新しました");
      } else {
        const res = await POST(`/api/training-records`, { body });
        if (!res.ok) return error("保存に失敗しました");
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
          {/* 年月と日付選択を統合 */}
          <div className="bg-white rounded-lg shadow mb-6 p-6">
            <div className="flex items-center justify-center gap-4 mb-6">
              <input readOnly type="number" value={Year} className="w-24 text-center border rounded-md px-3 py-2" />
              <span className="text-lg font-medium">年</span>
              <input readOnly type="number" value={Month} className="w-20 text-center border rounded-md px-3 py-2" />
              <span className="text-lg font-medium">月</span>
            </div>
            <MonthlyStrip year={Year} month={Month} selectedDay={selectedDay} onSelectDay={setSelectedDay} />
          </div>
          {/* トレーニング詳細 */}
          <div className="bg-white rounded-lg shadow mb-6 p-6">
            <div className="flex items-start justify-between gap-4">
              <TrainingSessionMetaEditor
                value={{
                  startedAt: form.getValues("startedAt"),
                  endedAt: form.getValues("endedAt"),
                  place: form.getValues("place"),
                  conditionLevel: form.getValues("conditionLevel"),
                  trainingPart: form.getValues("trainingPart"),
                }}
                onChange={(next) => {
                  form.setValue("startedAt", next.startedAt);
                  form.setValue("endedAt", next.endedAt);
                  form.setValue("place", next.place);
                  form.setValue("conditionLevel", next.conditionLevel);
                  form.setValue("trainingPart", next.trainingPart);
                }}
              />
              <button onClick={form.handleSubmit(handleSubmit)} className="px-6 py-2 rounded-md bg-booking-600 text-white disabled:opacity-50 whitespace-nowrap self-center">
                {isUpdate ? "修正" : "登録"}
              </button>
            </div>

            {/* 区切り線 */}
            <div className="my-6 border-t border-gray-200"></div>

            {/* 種目＋セット */}
            <TrainingExercisesEditor rows={rows ?? defaultValues.exercises} onChangeRows={(next) => form.setValue("exercises", next, { shouldDirty: true })} />
          </div>
        </div>
      </div>
    </FormProvider>
  );
};

export default TrainingRecordEditor;
