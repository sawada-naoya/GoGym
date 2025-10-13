"use client";
import React from "react";
import MonthlyStrip from "./MonthlyStrip";
import SessionMetaForm from "./SessionMetaForm";
import TrainingTable from "./TrainingTables";
import type { Exercise, SessionMeta } from "./types";

const HomeClient = ({ initialYear, initialMonth, initialDay, initialExercises }: { initialYear: number; initialMonth: number; initialDay: number; initialExercises: Exercise[] }) => {
  const [year, setYear] = React.useState(initialYear);
  const [month, setMonth] = React.useState(initialMonth);
  const [selectedDay, setSelectedDay] = React.useState(initialDay);

  const [meta, setMeta] = React.useState<SessionMeta>({
    startedAt: null,
    endedAt: null,
    place: "",
    muscle: "",
    condition: null,
  });

  // テーブル用のダミー1行（MVP）
  const [rows, setRows] = React.useState([{ exerciseId: "", sets: Array.from({ length: 5 }, (_, i) => ({ setNumber: i + 1, weight: "", rep: "" })) }]);

  const toISO = (hmm?: string | null) => {
    if (!hmm) return null;
    const [hh, mm] = hmm.split(":").map(Number);
    const d = new Date(year, month - 1, selectedDay, hh, mm, 0);
    return d.toISOString();
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="container mx-auto px-4 max-w-7xl">
        {/* 年月表示（必要ならナビボタンを後で） */}
        <div className="bg-white rounded-lg shadow mb-6 p-6">
          <div className="flex items-center justify-center gap-4">
            <input readOnly type="number" value={year} className="w-24 px-3 py-2 border border-gray-300 rounded-md text-center" />
            <span className="text-lg font-medium">年</span>
            <input readOnly type="number" value={month} className="w-20 px-3 py-2 border border-gray-300 rounded-md text-center" />
            <span className="text-lg font-medium">月</span>
          </div>
        </div>

        {/* カレンダー（日選択） */}
        <MonthlyStrip year={year} month={month} selectedDay={selectedDay} onSelectDay={setSelectedDay} />

        {/* 時間・場所・部位・コンディション */}
        <SessionMetaForm value={meta} onChange={(next) => setMeta(next)} />

        {/* トレーニング記録テーブル */}
        <TrainingTable exercises={initialExercises} rows={rows} onChangeRows={setRows} />

        {/* 保存ペイロード例（APIできたらPOST） */}
        {/* <button onClick={handleSave}>保存</button> */}
      </div>
    </div>
  );
};

export default HomeClient;
