"use client";
import React from "react";
import type { Exercise } from "./types";

type Muscle = "" | "腕" | "胸" | "脚" | "肩" | "背中" | "体幹";

const MUSCLE_CHIPS: { label: string; value: Muscle }[] = [
  { label: "すべて", value: "" },
  { label: "腕", value: "腕" },
  { label: "胸", value: "胸" },
  { label: "脚", value: "脚" },
  { label: "肩", value: "肩" },
  { label: "背中", value: "背中" },
  { label: "体幹", value: "体幹" },
];

const norm = (s: string) => s.normalize("NFKC").toLowerCase();

export default function TrainingTable({ exercises, rows, onChangeRows }: { exercises: Exercise[]; rows: { exerciseId: string; sets: { setNumber: number; weight: number | string; rep: number | string }[] }[]; onChangeRows: (rows: any) => void }) {
  // --- ソフトフィルタ（部位）状態 ---
  const [muscleFilter, setMuscleFilter] = React.useState<Muscle>("");

  // --- モック: 最近/お気に入り/Big3（本番はAPI/DBから） ---
  const recentIds = React.useMemo(() => ["bench", "squat"], []);
  const favoriteIds = React.useMemo(() => ["lat-pd"], []);
  const big3Names = new Set(["ベンチプレス", "スクワット", "デッドリフト"]);

  // exercises に muscle 情報が無い前提でも動くよう、簡易マッピング
  // 実運用では Exercise に muscle を持たせると良い（部位候補の優先度付けに使用）
  const withMeta = React.useMemo(() => {
    const guess = (name: string): Muscle => {
      if (/ベンチ|チェスト|胸/i.test(name)) return "胸";
      if (/スクワット|脚|レッグ/i.test(name)) return "脚";
      if (/ショルダー|肩/i.test(name)) return "肩";
      if (/ラット|背中|ロー|デッド/i.test(name)) return "背中";
      if (/カール|トライセプ|アーム|腕/i.test(name)) return "腕";
      if (/体幹|プランク|コア/i.test(name)) return "体幹";
      return "";
    };
    return exercises.map((e) => ({
      ...e,
      muscle: guess(e.name) as Muscle,
      isFavorite: favoriteIds.includes(e.id),
      isRecent: recentIds.includes(e.id),
      isBig3: big3Names.has(e.name),
    }));
  }, [exercises, favoriteIds, recentIds]);

  const updateCell = (rowIdx: number, setIdx: number, key: "weight" | "rep", val: string) => {
    const next = [...rows];
    next[rowIdx].sets[setIdx][key] = val;
    onChangeRows(next);
  };

  const changeExercise = (rowIdx: number, id: string) => {
    const next = [...rows];
    next[rowIdx].exerciseId = id;
    onChangeRows(next);
  };

  return (
    <div className="bg-white rounded-lg shadow overflow-x-auto">
      {/* ソフトフィルタ: 部位チップ列（選んでも“優先表示”するだけ。入力中は無視） */}
      <div className="px-4 pt-4 flex flex-wrap gap-2">
        {MUSCLE_CHIPS.map((chip) => (
          <button key={chip.label} type="button" onClick={() => setMuscleFilter(chip.value)} className={`px-3 py-1 rounded-full border text-sm ${muscleFilter === chip.value ? "bg-booking-600 text-white border-booking-600" : "border-gray-300 text-gray-700 hover:bg-gray-100"}`}>
            {chip.label}
          </button>
        ))}
      </div>

      <table className="w-full border-collapse mt-2">
        <thead>
          <tr className="border-b-2 border-gray-300">
            <th className="px-4 py-3 text-left font-medium text-gray-700 border-r border-gray-300 min-w-[220px]">種目</th>
            {[1, 2, 3, 4, 5].map((set) => (
              <th key={set} colSpan={2} className="px-4 py-3 text-center font-medium text-gray-700 border-r border-gray-300">
                {set}セット
              </th>
            ))}
          </tr>
          <tr className="border-b border-gray-300 bg-gray-50">
            <th className="px-4 py-2 border-r border-gray-300"></th>
            {[1, 2, 3, 4, 5].map((set) => (
              <React.Fragment key={set}>
                <th className="px-2 py-2 text-center text-sm font-medium text-gray-600 border-r border-gray-200">重量</th>
                <th className="px-2 py-2 text-center text-sm font-medium text-gray-600 border-r border-gray-300">回数</th>
              </React.Fragment>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, rowIdx) => (
            <React.Fragment key={rowIdx}>
              <tr className="hover:bg-gray-50">
                <td className="px-4 py-3 border-r border-gray-300 align-top border-b border-gray-300">
                  <ExerciseCombobox
                    all={withMeta}
                    muscleFilter={muscleFilter}
                    valueId={row.exerciseId}
                    onPick={(exId) => changeExercise(rowIdx, exId)}
                    onCreateNew={(name) => {
                      // 実APIできたら POST /exercises -> 返却ID に置換
                      const tmpId = crypto.randomUUID();
                      changeExercise(rowIdx, tmpId);
                    }}
                  />
                </td>

                {row.sets.map((s, i) => (
                  <React.Fragment key={i}>
                    <td className="px-2 py-2 border-r border-gray-200 border-b border-gray-200">
                      <input type="number" value={s.weight} placeholder="" onChange={(e) => updateCell(rowIdx, i, "weight", e.target.value)} className="w-full px-2 py-1 border border-gray-300 rounded text-center focus:outline-none focus:ring-1 focus:ring-booking-500" />
                    </td>
                    <td className="px-2 py-2 border-r border-gray-300 border-b border-gray-200">
                      <input type="number" value={s.rep} placeholder="" onChange={(e) => updateCell(rowIdx, i, "rep", e.target.value)} className="w-full px-2 py-1 border border-gray-300 rounded text-center focus:outline-none focus:ring-1 focus:ring-booking-500" />
                    </td>
                  </React.Fragment>
                ))}
              </tr>

              <tr className="border-b border-gray-300 hover:bg-gray-50">
                <td colSpan={11} className="px-4 py-2">
                  <input type="text" placeholder="メモ（任意）" className="w-full px-2 py-1 text-xs text-gray-600 border-0 focus:outline-none" />
                </td>
              </tr>
            </React.Fragment>
          ))}

          <tr className="h-32 border-b border-gray-200">
            <td colSpan={12} className="px-4 py-3 text-center text-gray-400">
              <button
                className="text-booking-600 hover:text-booking-700 font-medium"
                onClick={() => {
                  onChangeRows([...rows, { exerciseId: "", sets: Array.from({ length: 5 }, (_, i) => ({ setNumber: i + 1, weight: "", rep: "" })) }]);
                }}
              >
                + 種目を追加
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}

/** 検索コンボボックス（部位チップは“優先表示”にのみ使用。入力中は全検索） */
function ExerciseCombobox({ all, muscleFilter, valueId, onPick, onCreateNew }: { all: (Exercise & { muscle: Muscle; isFavorite?: boolean; isRecent?: boolean; isBig3?: boolean })[]; muscleFilter: Muscle; valueId: string; onPick: (id: string) => void; onCreateNew: (name: string) => void }) {
  const [open, setOpen] = React.useState(false);
  const [q, setQ] = React.useState("");
  const selected = all.find((e) => e.id === valueId) || null;

  // 並び順ロジック：
  // 1) 入力中: 全体から名前ヒット（部位無視）
  // 2) 未入力: ソフトグループ（最近→お気に入り→Big3）を先頭に、muscleFilter一致を上に優先
  const filtered = React.useMemo(() => {
    if (q) {
      const hits = all.filter((e) => norm(e.name).includes(norm(q)));
      return hits.slice(0, 20);
    }
    // 未入力時の優先ソート
    const score = (e: any) => (e.isRecent ? 100 : 0) + (e.isFavorite ? 50 : 0) + (e.isBig3 ? 25 : 0) + (muscleFilter && e.muscle === muscleFilter ? 10 : 0);
    return [...all].sort((a, b) => score(b) - score(a)).slice(0, 30);
  }, [all, q, muscleFilter]);

  return (
    <div className="relative">
      <input
        value={q}
        onChange={(e) => {
          setQ(e.target.value);
          setOpen(true);
        }}
        onFocus={() => setOpen(true)}
        placeholder={selected?.name || "種目を検索 / 追加"}
        className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-booking-500"
      />
      {open && (
        <div className="absolute z-20 mt-1 w-80 max-w-[90vw] bg-white border border-gray-200 rounded shadow">
          {q && filtered.length === 0 && (
            <button
              className="w-full text-left px-3 py-2 hover:bg-gray-50"
              onMouseDown={(e) => e.preventDefault()}
              onClick={() => {
                onCreateNew(q);
                setQ("");
                setOpen(false);
              }}
            >
              「{q}」を新規追加
            </button>
          )}

          {!q && <div className="px-3 pt-2 pb-1 text-[11px] text-gray-500">最近 / お気に入り / Big3 を優先表示</div>}

          <ul className="max-h-64 overflow-auto">
            {filtered.map((e) => (
              <li key={e.id}>
                <button
                  className="w-full flex items-center justify-between px-3 py-2 hover:bg-gray-50"
                  onMouseDown={(ev) => ev.preventDefault()}
                  onClick={() => {
                    onPick(e.id);
                    setQ("");
                    setOpen(false);
                  }}
                >
                  <span>{e.name}</span>
                  <span className="text-[10px] text-gray-400">
                    {e.isRecent ? "最近 " : ""}
                    {e.isFavorite ? "★ " : ""}
                    {e.isBig3 ? "Big3 " : ""}
                    {e.muscle || ""}
                  </span>
                </button>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}
