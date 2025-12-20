import WorkoutContent from "./content";
import { extractDateParts } from "./_lib/utils";
import { buildEmptyDTO } from "./_lib/types";
import { SeedWorkoutParts } from "@/app/api/workouts/seed/route";
import { GetWorkoutRecords } from "@/app/api/workouts/records/route";
import { GetWorkoutParts } from "@/app/api/workouts/parts/route";

export const dynamic = "force-dynamic";

type Props = {
  params: { user_id: string };
  searchParams?: { date?: string };
};

const Page = async ({ searchParams }: Props) => {
  // 部位データをシード（初回のみ作成、idempotent）
  await SeedWorkoutParts();

  // SSRで並列取得
  const [dtoResponse, partsResponse] = await Promise.all([
    GetWorkoutRecords({ date: searchParams?.date }),
    GetWorkoutParts(),
  ]);

  // NextResponseからJSONデータを取得
  const dto =
    dtoResponse.status === 200 ? await dtoResponse.json() : buildEmptyDTO();
  const parts = partsResponse.status === 200 ? await partsResponse.json() : [];

  // バックエンドから返された日付を使用（通常は必ず返される）
  const { year, month, day } = extractDateParts(dto.performed_date);

  return (
    <WorkoutContent
      Year={year}
      Month={month}
      Day={day}
      defaultValues={dto}
      availableParts={parts}
      isUpdate={!!dto.id}
    />
  );
};

export default Page;
