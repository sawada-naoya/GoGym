"use client";

import { WorkoutProvider } from "../_providers/WorkoutProvider";
import { WorkoutView } from "../_components/WorkoutView";
import type { WorkoutFormDTO, WorkoutPartDTO } from "@/types/workout";

/**
 * WorkoutContainer
 *
 * 責務:
 * - 初期データ（defaultValues, availableParts）を受け取る
 * - WorkoutProviderでラップ
 * - WorkoutViewを呼び出す
 *
 * ロジックは一切持たない（全てProviderに委譲）
 */

type Props = {
  defaultValues: WorkoutFormDTO;
  availableParts: WorkoutPartDTO[];
};

export default function WorkoutContainer({
  defaultValues,
  availableParts,
}: Props) {
  return (
    <WorkoutProvider
      defaultValues={defaultValues}
      availableParts={availableParts}
    >
      <WorkoutView />
    </WorkoutProvider>
  );
}
