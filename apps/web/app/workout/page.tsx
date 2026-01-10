import WorkoutContent from "./content";
import { extractDateParts } from "@/features/workout/lib/utils";
import {
  buildEmptyDTO,
  type WorkoutFormDTO,
  type WorkoutPartDTO,
} from "@/types/workout";
import { getServerAccessToken } from "@/features/auth/server";

export const dynamic = "force-dynamic";

const API_BASE = process.env.NEXT_PUBLIC_API_URL;

type Props = {
  params: { user_id: string };
  searchParams?: { date?: string };
};

const Page = async ({ searchParams }: Props) => {
  const sp = await searchParams;
  const token = await getServerAccessToken();

  if (!token || !API_BASE) {
    return <div>認証エラー</div>;
  }

  // 部位データをシード（初回のみ作成、idempotent）
  await fetch(`${API_BASE}/api/v1/workouts/seed`, {
    method: "POST",
    headers: { Authorization: `Bearer ${token}` },
    cache: "no-store",
  });

  // SSRで並列取得
  const query = new URLSearchParams();
  if (sp?.date) query.set("date", sp.date);
  const queryString = query.toString() ? `?${query.toString()}` : "";

  const [recordsRes, partsRes] = await Promise.all([
    fetch(`${API_BASE}/api/v1/workouts/records${queryString}`, {
      headers: { Authorization: `Bearer ${token}` },
      cache: "no-store",
    }),
    fetch(`${API_BASE}/api/v1/workouts/parts`, {
      headers: { Authorization: `Bearer ${token}` },
      cache: "no-store",
    }),
  ]);

  const dto: WorkoutFormDTO = recordsRes.ok
    ? await recordsRes.json()
    : buildEmptyDTO();
  const parts: WorkoutPartDTO[] = partsRes.ok ? await partsRes.json() : [];

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
