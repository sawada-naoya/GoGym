"use client";

import { useEffect } from "react";
import { useForm, UseFormReturn } from "react-hook-form";
import type { WorkoutFormDTO } from "@/types/workout";

type UseWorkoutFormProps = {
  defaultValues: WorkoutFormDTO;
  onSubmit: (data: WorkoutFormDTO) => Promise<void>;
};

type UseWorkoutFormReturn = {
  form: UseFormReturn<WorkoutFormDTO>;
  handleSubmit: () => Promise<void>;
  isSubmitting: boolean;
};

/**
 * ワークアウトフォーム管理hook
 *
 * 責務:
 * - React Hook Formの初期化
 * - defaultValuesの変更検知とリセット
 * - submit処理の統合
 */
export const useWorkoutForm = ({
  defaultValues,
  onSubmit,
}: UseWorkoutFormProps): UseWorkoutFormReturn => {
  const form = useForm<WorkoutFormDTO>({
    defaultValues,
    mode: "onBlur",
  });

  // defaultValuesが変わったらフォームをリセット（1箇所のみ）
  useEffect(() => {
    form.reset(defaultValues);
  }, [defaultValues, form]);

  const handleSubmit = async () => {
    await form.handleSubmit(onSubmit)();
  };

  return {
    form,
    handleSubmit,
    isSubmitting: form.formState.isSubmitting,
  };
};
