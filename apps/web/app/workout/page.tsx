import WorkoutContent from "./content";
import { extractDateParts } from "./_lib/utils";
import type { WorkoutFormDTO, WorkoutPartDTO } from "./_lib/types";
import { buildEmptyDTO } from "./_lib/types";
import { getServerAccessToken } from "@/lib/auth-helpers";
import {
  seedWorkoutParts,
  getWorkoutParts,
  getWorkoutRecord,
} from "@/lib/bff/workout";

export const dynamic = "force-dynamic";

type Props = {
  params: { user_id: string };
  searchParams?: { date?: string };
};

const Page = async ({ searchParams }: Props) => {
  const token = await getServerAccessToken();
  if (!token) {
    // ログインしていない場合はログイン画面へリダイレクト
    return (
      <div className="p-4">
        <p className="text-red-500">
          ログインが必要です。
          <a href="/auth/login" className="underline">
            ログインページへ
          </a>
        </p>
      </div>
    );
  }

  // 部位データをシード（初回のみ作成、idempotent）
  await seedWorkoutParts();

  // SSRで並列取得
  const [dto, parts] = await Promise.all([
    getWorkoutRecord({ date: searchParams?.date }).catch(() => buildEmptyDTO()),
    getWorkoutParts(),
  ]);

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
