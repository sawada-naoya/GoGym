"use client";
import { useEffect, useState } from "react";
import { useForm, FormProvider, useWatch } from "react-hook-form";
import WorkoutMetadataEditor from "./_components/WorkoutMetadataEditor";
import WorkoutExercisesEditor from "./_components/WorkoutExercisesEditor";
import { WorkoutFormDTO, WorkoutPartDTO, ExerciseDTO } from "./_lib/types";
import { transformFormDataForSubmit } from "./_lib/utils";
import { buildEmptyDTO } from "./_lib/types";
import { useBanner } from "@/components/Banner";
import { formatDate } from "@/lib/time";
import { GET, POST, PUT } from "@/lib/api";

type Props = {
  Year: number;
  Month: number;
  Day: number;
  defaultValues: WorkoutFormDTO;
  availableParts: WorkoutPartDTO[];
  isUpdate: boolean;
  token: string;
};

/**
 * ワークアウト記録を取得
 */
const fetchWorkoutRecord = async (token: string, date?: string, partID?: number | null): Promise<WorkoutFormDTO> => {
  const params = new URLSearchParams();
  if (date) params.append("date", date);
  if (partID) params.append("part_id", partID.toString());

  const queryString = params.toString();
  const url = queryString ? `/api/v1/workouts/records?${queryString}` : "/api/v1/workouts/records";
  const res = await GET(url, { token: token ?? undefined });

  if (!res.ok) {
    return buildEmptyDTO();
  }

  const display = res.data as WorkoutFormDTO;

  if (!display || !display.id) {
    const emptyExercises = Array.from({ length: 1 }, () => ({
      id: null,
      name: "",
      workout_part_id: null,
      sets: Array.from({ length: 1 }, (_, i) => ({
        id: null,
        set_number: i + 1,
        weight_kg: 0,
        reps: 0,
        note: null,
      })),
    }));

    return {
      id: null,
      performed_date: display?.performed_date || "",
      started_at: null,
      ended_at: null,
      place: display?.place || "",
      note: null,
      condition_level: null,
      workout_part: { id: null, name: null, source: null },
      exercises: display?.exercises?.length > 0 ? display.exercises : emptyExercises,
    };
  }

  return {
    id: display.id,
    performed_date: display.performed_date,
    started_at: display.started_at ?? null,
    ended_at: display.ended_at ?? null,
    place: display.place ?? "",
    note: display.note ?? null,
    condition_level: display.condition_level ?? null,
    workout_part: {
      id: display.workout_part?.id ?? null,
      name: display.workout_part?.name ?? null,
      source: display.workout_part?.source ?? null,
    },
    exercises: display.exercises.map((ex) => ({
      id: ex.id,
      name: ex.name,
      workout_part_id: ex.workout_part_id,
      sets: (ex.sets ?? []).map((s) => ({
        id: s.id,
        set_number: s.set_number,
        weight_kg: s.weight_kg,
        reps: s.reps,
        note: s.note ?? null,
      })),
    })),
  };
};

/**
 * ワークアウト部位リストを取得
 */
const fetchWorkoutParts = async (token: string): Promise<WorkoutPartDTO[]> => {
  const res = await GET<WorkoutPartDTO[]>("/api/v1/workouts/parts", { token: token ?? undefined });

  if (!res.ok || !Array.isArray(res.data)) {
    return [];
  }

  return res.data;
};

/**
 * ワークアウト記録を作成
 */
async function createWorkoutRecord(token: string, body: any) {
  if (!token) {
    return { ok: false, error: "認証エラー" };
  }

  const res = await POST("/api/v1/workouts/records", { body, token: token });
  return { ok: res.ok, error: res.ok ? null : "保存に失敗しました" };
}

/**
 * ワークアウト記録を更新
 */
async function updateWorkoutRecord(token: string, id: number, body: any) {
  if (!token) {
    return { ok: false, error: "認証エラー" };
  }

  const res = await PUT(`/api/v1/workouts/records/${id}`, { body, token: token });
  return { ok: res.ok, error: res.ok ? null : "更新に失敗しました" };
}

/**
 * 前回のエクササイズ記録を取得
 */
async function fetchLastExerciseRecord(token: string, exerciseID: number): Promise<ExerciseDTO | null> {
  if (!token || !exerciseID) {
    return null;
  }

  const res = await GET<ExerciseDTO>(`/api/v1/workouts/exercises/${exerciseID}/last`, { token });

  if (!res.ok || !res.data) {
    return null;
  }

  return res.data;
}

const WorkoutContent = ({ Year, Month, Day, defaultValues, availableParts: initialParts, isUpdate, token }: Props) => {
  const { success, error } = useBanner();
  const [selectedDay, setSelectedDay] = useState(Day);
  const [selectedYear, setSelectedYear] = useState(Year);
  const [selectedMonth, setSelectedMonth] = useState(Month);
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
      try {
        const date = formatDate(selectedYear, selectedMonth, selectedDay);
        const data = await fetchWorkoutRecord(token, date, selectedPart.id);
        // 部位とexercisesのみ更新（メタデータは維持）
        form.setValue("exercises", data.exercises, { shouldDirty: false });
      } catch (err) {
        console.error("Failed to fetch part records:", err);
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
            dataKey={`${selectedYear}-${selectedMonth}-${selectedDay}-${selectedPart?.id || "none"}`}
            onFetchLastRecord={fetchLastExerciseRecord}
            token={token}
          />
        </div>
      </div>
    </FormProvider>
  );
};

export default WorkoutContent;
