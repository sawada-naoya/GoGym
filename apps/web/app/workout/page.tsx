import { GET } from "@/lib/api";
import { getServerSession } from "next-auth";
import { redirect } from "next/navigation";
import WorkoutRecordEditor from "./_components/WorkoutRecordEditor";
import { WorkoutFormDTO } from "./_lib/types";

export const dynamic = "force-dynamic";

const toHHmm = (iso: string | null): string | null => {
  if (!iso) return null;
  const d = new Date(iso);
  const hh = String(d.getHours()).padStart(2, "0");
  const mm = String(d.getMinutes()).padStart(2, "0");
  return `${hh}:${mm}`;
};

const buildEmptyDTO = (): WorkoutFormDTO => ({
  id: null,
  performed_date: "",
  started_at: null,
  ended_at: null,
  place: "",
  note: null,
  condition_level: null,
  workout_part: { id: null, name: null, source: null },
  exercises: [
    {
      id: null,
      name: "",
      workout_part_id: null,
      is_default: 0,
      sets: Array.from({ length: 5 }, (_, i) => ({
        id: null,
        set_number: i + 1,
        weight_kg: "",
        reps: "",
        note: null,
      })),
    },
  ],
});

const fetchWorkoutRecord = async (date?: string): Promise<WorkoutFormDTO> => {
  // dateがある場合はクエリパラメータに含める、ない場合はバックエンドが今日のJST日付を使用
  const url = date ? `/api/v1/workouts/records?date=${date}` : "/api/v1/workouts/records";
  const res = await GET(url);

  if (!res.ok) {
    // エラーの場合は空のDTOを返す
    return buildEmptyDTO();
  }

  const display = res.data as WorkoutFormDTO;

  // バックエンドから返されたデータをそのまま使用（日付も含まれている）
  if (!display || !display.id) {
    // レコードがない場合もバックエンドから日付は返されている
    return {
      ...buildEmptyDTO(),
      performed_date: display.performed_date || "",
      place: display.place || "",
      exercises: display.exercises || buildEmptyDTO().exercises,
    };
  }

  return {
    id: display.id,
    performed_date: display.performed_date,
    started_at: toHHmm(display.started_at),
    ended_at: toHHmm(display.ended_at),
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
      is_default: ex.is_default,
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

type Props = {
  params: { user_id: string };
  searchParams?: { date?: string };
};

const Page = async ({ searchParams }: Props) => {
  const session = await getServerSession();

  // 未ログインの場合はトップページにリダイレクト
  if (!session?.user) {
    redirect("/");
  }

  // バックエンドに日付処理を任せる（dateがundefinedの場合、バックエンドが今日のJST日付を使用）
  const dto = await fetchWorkoutRecord(searchParams?.date);

  // バックエンドから返された日付を使用
  const date = dto.performed_date;
  const year = Number(date.slice(0, 4));
  const month = Number(date.slice(5, 7));
  const day = Number(date.slice(8, 10));

  return <WorkoutRecordEditor Year={year} Month={month} Day={day} defaultValues={dto} isUpdate={!!dto.id} />;
};

export default Page;
