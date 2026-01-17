import WorkoutContainer from "./_containers/WorkoutContainer";
import type { WorkoutFormDTO, WorkoutPartDTO } from "@/types/workout";
import {
  seedWorkoutParts,
  getWorkoutRecords,
  getWorkoutParts,
} from "@/features/workout/actions";
import { buildEmptyDTO } from "@/features/workout/lib/transforms";

export const dynamic = "force-dynamic";

type Props = {
  params: { user_id: string };
  searchParams?: { date?: string };
};

const Page = async ({ searchParams }: Props) => {
  const sp = await searchParams;

  // 部位データをシード（初回のみ作成、idempotent）
  await seedWorkoutParts();

  // SSRで並列取得
  const dateStr = sp?.date;
  const [recordsResult, partsResult] = await Promise.all([
    getWorkoutRecords(dateStr ? { date: dateStr } : undefined),
    getWorkoutParts(),
  ]);
  const dto: WorkoutFormDTO =
    recordsResult.success && recordsResult.data
      ? recordsResult.data
      : buildEmptyDTO();
  const parts: WorkoutPartDTO[] =
    partsResult.success && partsResult.data ? partsResult.data : [];

  return <WorkoutContainer defaultValues={dto} availableParts={parts} />;
};

export default Page;
