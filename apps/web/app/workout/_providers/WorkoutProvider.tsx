"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  type ReactNode,
} from "react";
import { FormProvider, useWatch, type UseFormReturn } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { useBanner } from "@/components/Banner";
import { useWorkoutDate } from "@/features/workout/hooks/useWorkoutDate";
import { useWorkoutForm } from "@/features/workout/hooks/useWorkoutForm";
import { useLastRecords } from "@/features/workout/hooks/useLastRecords";
import { transformFormDataForSubmit } from "@/features/workout/lib/transforms";
import {
  createWorkoutRecord,
  updateWorkoutRecord,
  getWorkoutParts,
} from "@/features/workout/actions";
import type {
  WorkoutFormDTO,
  WorkoutPartDTO,
  ExerciseDTO,
} from "@/types/workout";

type WorkoutContextValue = {
  // 日付管理
  year: number;
  month: number;
  day: number;
  dateStr: string;
  setYear: (year: number) => void;
  setMonth: (month: number) => void;
  setDay: (day: number) => void;
  setDate: (year: number, month: number, day: number) => void;

  // フォーム管理
  form: UseFormReturn<WorkoutFormDTO>;
  handleSubmit: () => Promise<void>;
  isSubmitting: boolean;

  // 種目の前回記録
  fetchLastRecord: (exerciseId: number) => Promise<ExerciseDTO | null>;

  // UI状態
  availableParts: WorkoutPartDTO[];
  setAvailableParts: (parts: WorkoutPartDTO[]) => void;
  selectedPartId: number | null;
  setSelectedPartId: (id: number | null) => void;
  refetchWorkoutParts: () => Promise<void>;

  // その他
  isUpdate: boolean;
  allExercises: WorkoutFormDTO["exercises"];
};

const WorkoutContext = createContext<WorkoutContextValue | null>(null);

/**
 * WorkoutContextを使用するカスタムhook
 * コンポーネント内で const { year, month, ... } = useWorkoutContext(); として使用
 */
export const useWorkoutContext = () => {
  const context = useContext(WorkoutContext);
  if (!context) {
    throw new Error("useWorkoutContext must be used within WorkoutProvider");
  }
  return context;
};

type WorkoutProviderProps = {
  defaultValues: WorkoutFormDTO;
  availableParts: WorkoutPartDTO[];
  children: ReactNode;
};

export const WorkoutProvider = ({
  defaultValues,
  availableParts: initialParts,
  children,
}: WorkoutProviderProps) => {
  const { t } = useTranslation("common");
  const { success, error } = useBanner();

  // 日付管理
  const { year, month, day, dateStr, setYear, setMonth, setDay, setDate } =
    useWorkoutDate();

  // 前回記録管理
  const { fetchLastRecord } = useLastRecords();

  // フォーム管理
  const {
    form,
    handleSubmit: handleFormSubmit,
    isSubmitting,
  } = useWorkoutForm({
    defaultValues,
    onSubmit: async (data) => {
      const body = transformFormDataForSubmit(data, year, month, day) as any;

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

  // ==================== UI State ====================
  const [availableParts, setAvailableParts] =
    useState<WorkoutPartDTO[]>(initialParts);

  const [selectedPartId, setSelectedPartId] = useState<number | null>(() => {
    return defaultValues.exercises?.[0]?.workout_part_id ?? null;
  });

  // defaultValuesが変わったらselectedPartIdも更新
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

  const isUpdate = !!defaultValues.id;

  const value: WorkoutContextValue = {
    // 日付管理
    year,
    month,
    day,
    dateStr,
    setYear,
    setMonth,
    setDay,
    setDate,

    // フォーム管理
    form,
    handleSubmit: handleFormSubmit,
    isSubmitting,

    // 種目の前回記録
    fetchLastRecord,

    // UI状態
    availableParts,
    setAvailableParts,
    selectedPartId,
    setSelectedPartId,
    refetchWorkoutParts,

    // その他
    isUpdate,
    allExercises: allExercises ?? defaultValues.exercises,
  };

  return (
    <FormProvider {...form}>
      <WorkoutContext.Provider value={value}>
        {children}
      </WorkoutContext.Provider>
    </FormProvider>
  );
};
