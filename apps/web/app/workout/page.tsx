import { auth } from "@/app/api/auth/[...nextauth]/authOptions";
import WorkoutContent from "./content";
import { extractDateParts } from "./_lib/utils";
import { GET, POST } from "@/lib/api";
import type { WorkoutFormDTO, WorkoutPartDTO } from "./_lib/types";
import { buildEmptyDTO } from "./_lib/types";

export const dynamic = "force-dynamic";

type Props = {
  params: { user_id: string };
  searchParams?: { date?: string };
};

/**
 * ワークアウト部位をシード（初期データ作成）
 */
const seedWorkoutParts = async (token: string): Promise<void> => {
  await POST("/api/v1/workouts/seed", { token: token ?? undefined });
};

/**
 * ワークアウト部位リストを取得
 */
const fetchWorkoutParts = async (token: string): Promise<WorkoutPartDTO[]> => {
  const res = await GET<WorkoutPartDTO[]>("/api/v1/workouts/parts", { token: token ?? undefined });

  if (!res.ok || !Array.isArray(res.data)) {
    return [];
  }

  return res.data;
};

/**
 * ワークアウト記録を取得
 */
const fetchWorkoutRecord = async (token: string, date?: string, partID?: number | null): Promise<WorkoutFormDTO> => {
  const params = new URLSearchParams();
  if (date) params.append("date", date);
  if (partID) params.append("part_id", partID.toString());

  const queryString = params.toString();
  const url = queryString ? `/api/v1/workouts/records?${queryString}` : "/api/v1/workouts/records";
  const res = await GET(url, { token: token ?? undefined });

  if (!res.ok) {
    return buildEmptyDTO();
  }

  const display = res.data as WorkoutFormDTO;

  if (!display || !display.id) {
    const emptyExercises = Array.from({ length: 1 }, () => ({
      id: null,
      name: "",
      workout_part_id: null,
      sets: Array.from({ length: 1 }, (_, i) => ({
        id: null,
        set_number: i + 1,
        weight_kg: "",
        reps: "",
        note: null,
      })),
    }));

    return {
      id: null,
      performed_date: display?.performed_date || "",
      started_at: null,
      ended_at: null,
      place: display?.place || "",
      note: null,
      condition_level: null,
      workout_part: { id: null, name: null, source: null },
      exercises: display?.exercises?.length > 0 ? display.exercises : emptyExercises,
    };
  }

  return {
    id: display.id,
    performed_date: display.performed_date,
    started_at: display.started_at ?? null,
    ended_at: display.ended_at ?? null,
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

const Page = async ({ searchParams }: Props) => {
  const session = await auth();
  const token = session?.user?.accessToken!;

  // 部位データをシード（初回のみ作成、idempotent）
  await seedWorkoutParts(token);

  // SSRで並列取得
  const [dto, parts] = await Promise.all([fetchWorkoutRecord(token, searchParams?.date), fetchWorkoutParts(token)]);

  // バックエンドから返された日付を使用（通常は必ず返される）
  const { year, month, day } = extractDateParts(dto.performed_date);

  return <WorkoutContent Year={year} Month={month} Day={day} defaultValues={dto} availableParts={parts} isUpdate={!!dto.id} token={token} />;
};

export default Page;
