import { GET } from "@/lib/api";
import TrainingRecordEditor from "./_components/TrainingRecordEditor";
import { TrainingFormDTO } from "./_components/types";

export const dynamic = "force-dynamic";

const toHHmm = (iso: string | null): string | null => {
  if (!iso) return null;
  const d = new Date(iso);
  const hh = String(d.getHours()).padStart(2, "0");
  const mm = String(d.getMinutes()).padStart(2, "0");
  return `${hh}:${mm}`;
};

const buildEmptyDTO = (date: string): TrainingFormDTO => ({
  id: null,
  performedDate: date,
  startedAt: null,
  endedAt: null,
  place: "",
  note: null,
  conditionLevel: null,
  trainingPart: { id: null, name: null, source: null },
  exercises: [
    {
      id: null,
      name: "",
      trainingPartId: null,
      isDefault: 0,
      sets: Array.from({ length: 5 }, (_, i) => ({
        id: null,
        setNumber: i + 1,
        weightKg: "",
        reps: "",
        note: null,
      })),
    },
  ],
});

const fetchDTO = async (date: string): Promise<TrainingFormDTO> => {
  const res = await GET(`/api/training-records?date=${date}`);
  if (!res.ok) {
    return buildEmptyDTO(date);
  }
  const display = res.data as TrainingFormDTO;
  if (!display || !display.id) {
    return buildEmptyDTO(date);
  }
  return {
    id: display.id,
    performedDate: display.performedDate,
    startedAt: toHHmm(display.startedAt),
    endedAt: toHHmm(display.endedAt),
    place: display.place ?? "",
    note: display.note ?? null,
    conditionLevel: display.conditionLevel ?? null,
    trainingPart: {
      id: display.trainingPart?.id ?? null,
      name: display.trainingPart?.name ?? null,
      source: display.trainingPart?.source ?? null,
    },
    exercises: display.exercises.map((ex) => ({
      id: ex.id,
      name: ex.name,
      trainingPartId: ex.trainingPartId,
      isDefault: ex.isDefault,
      sets: (ex.sets ?? []).map((s) => ({
        id: s.id,
        setNumber: s.setNumber,
        weightKg: s.weightKg, // 数値→フォームではそのまま扱える
        reps: s.reps,
        note: s.note,
      })),
    })),
  };
};

const Page = async ({ searchParams }: { searchParams?: { date?: string } }) => {
  const nowJST = new Date(Date.now() + 9 * 60 * 60 * 1000);
  const date = searchParams?.date ?? [nowJST.getUTCFullYear(), String(nowJST.getUTCMonth() + 1).padStart(2, "0"), String(nowJST.getUTCDate()).padStart(2, "0")].join("-");
  const dto = await fetchDTO(date);
  const year = Number(date.slice(0, 4));
  const month = Number(date.slice(5, 7));
  const day = Number(date.slice(8, 10));

  return (
    <TrainingRecordEditor
      Year={year}
      Month={month}
      Day={day}
      defaultValues={dto} // ← フォームにそのまま入る！
      isUpdate={!!dto.id} // ← 既存かどうかの判定もpage側で済ませる
    />
  );
};

export default Page;
