// Server: 初期データ取得だけ担当
export const dynamic = "force-dynamic";

type Exercise = { id: string; name: string; isBuiltin?: boolean };

async function fetchInitialExercises(): Promise<Exercise[]> {
  // TODO: 実APIができたら fetch(`${process.env.NEXT_PUBLIC_API}/exercises`)
  return [
    { id: "bench", name: "ベンチプレス", isBuiltin: true },
    { id: "squat", name: "スクワット", isBuiltin: true },
    { id: "deadlift", name: "デッドリフト", isBuiltin: true },
  ];
}

import HomeClient from "./_components/HomeClient";

export default async function Page() {
  const exercises = await fetchInitialExercises();
  const now = new Date();
  return <HomeClient initialYear={now.getFullYear()} initialMonth={now.getMonth() + 1} initialDay={now.getDate()} initialExercises={exercises} />;
}
