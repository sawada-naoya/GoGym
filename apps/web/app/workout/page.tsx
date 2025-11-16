import { auth } from "@/app/api/auth/[...nextauth]/authOptions";
import { redirect } from "next/navigation";
import WorkoutRecordEditor from "./_components/WorkoutRecordEditor";
import { fetchWorkoutRecord, fetchWorkoutParts, seedWorkoutParts } from "./_lib/api";
import { extractDateParts } from "./_lib/utils";

export const dynamic = "force-dynamic";

type Props = {
  params: { user_id: string };
  searchParams?: { date?: string };
};

const Page = async ({ searchParams }: Props) => {
  const session = await auth();
  const token = session?.user?.accessToken;

  // 未ログインの場合はトップページにリダイレクト
  if (!session?.user) {
    redirect("/");
  }

  // トークンがない場合はリダイレクト
  if (!token) {
    redirect("/");
  }

  // 部位データをシード（初回のみ作成、idempotent）
  await seedWorkoutParts(token);

  // SSRで並列取得
  const [dto, parts] = await Promise.all([fetchWorkoutRecord(token, searchParams?.date), fetchWorkoutParts(token)]);

  // バックエンドから返された日付を使用（通常は必ず返される）
  const { year, month, day } = extractDateParts(dto.performed_date);

  return <WorkoutRecordEditor Year={year} Month={month} Day={day} defaultValues={dto} availableParts={parts} isUpdate={!!dto.id} token={token} />;
};

export default Page;
