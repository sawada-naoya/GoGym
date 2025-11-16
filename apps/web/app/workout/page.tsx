import { getServerSession } from "next-auth";
import { authOptions } from "@/app/api/auth/[...nextauth]/route";
import { redirect } from "next/navigation";
import WorkoutRecordEditor from "./_components/WorkoutRecordEditor";
import { fetchWorkoutRecord, fetchWorkoutParts } from "./_lib/api";
import { extractDateParts } from "./_lib/utils";

export const dynamic = "force-dynamic";

type Props = {
  params: { user_id: string };
  searchParams?: { date?: string };
};

const Page = async ({ searchParams }: Props) => {
  const session = await getServerSession(authOptions);

  // 未ログインの場合はトップページにリダイレクト
  if (!session?.user) {
    redirect("/");
  }

  // SSRで並列取得
  const [dto, parts] = await Promise.all([fetchWorkoutRecord(searchParams?.date), fetchWorkoutParts()]);

  // バックエンドから返された日付を使用（通常は必ず返される）
  const { year, month, day } = extractDateParts(dto.performed_date);

  return <WorkoutRecordEditor Year={year} Month={month} Day={day} defaultValues={dto} availableParts={parts} isUpdate={!!dto.id} />;
};

export default Page;
