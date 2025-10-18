"use client";
import React from "react";
import type { TrainingFormDTO } from "./types";

type Row = TrainingFormDTO["exercises"][number];

type Props = {
  rows: Row[];
  onChangeRows: (rows: Row[]) => void;
};

const TrainingExercisesEditor: React.FC<Props> = ({ rows, onChangeRows }) => {
  const updateCell = (ri: number, si: number, key: "weightKg" | "reps", val: string) => {
    const next = structuredClone(rows) as Row[];
    (next[ri].sets[si] as any)[key] = val; // 入力中は "" を許容
    onChangeRows(next);
  };

  const updateSetNote = (ri: number, si: number, note: string) => {
    const next = structuredClone(rows) as Row[];
    next[ri].sets[si].note = note || null;
    onChangeRows(next);
  };

  const ensureFiveSets = (row: Row): Row => {
    const sets = row.sets ?? [];
    if (sets.length >= 5) return row;
    const baseLen = sets.length;
    const add: Row["sets"] = Array.from({ length: 5 - baseLen }, (_, i) => ({
      setNumber: baseLen + i + 1,
      weightKg: "" as const,
      reps: "" as const,
      note: null,
    }));
    return { ...row, sets: [...sets, ...add] };
  };

  const changeExerciseName = (ri: number, name: string) => {
    const next = structuredClone(rows) as Row[];
    next[ri].name = name;
    // 自由入力なので既存IDに紐づけない（新規扱い）。既存を維持したいならここで条件分岐。
    next[ri].id = name.trim() ? next[ri].id ?? null : null;
    next[ri] = ensureFiveSets(next[ri]);
    onChangeRows(next);
  };

  const addExerciseRow = () => {
    onChangeRows([
      ...rows,
      {
        id: null,
        name: "",
        trainingPartId: null,
        isDefault: 0,
        sets: Array.from({ length: 5 }, (_, i) => ({
          setNumber: i + 1,
          weightKg: "" as const,
          reps: "" as const,
          note: null,
        })),
      },
    ]);
  };

  return (
    <div className="bg-white rounded-lg shadow overflow-x-auto">
      <table className="w-full border-collapse">
        <thead>
          <tr className="border-b-2 border-gray-300">
            <th className="px-4 py-3 text-left font-medium text-gray-700 border-r border-gray-300 min-w-[240px]">種目</th>
            {[1, 2, 3, 4, 5].map((n) => (
              <th key={n} colSpan={2} className="px-4 py-3 text-center font-medium text-gray-700 border-r border-gray-300">
                {n}セット
              </th>
            ))}
          </tr>
          <tr className="border-b border-gray-300 bg-gray-50">
            <th className="px-4 py-2 border-r border-gray-300" />
            {[1, 2, 3, 4, 5].map((n) => (
              <React.Fragment key={n}>
                <th className="px-2 py-2 text-center text-sm font-medium text-gray-600 border-r border-gray-200">重量</th>
                <th className="px-2 py-2 text-center text-sm font-medium text-gray-600 border-r border-gray-300">回数</th>
              </React.Fragment>
            ))}
          </tr>
        </thead>

        <tbody>
          {rows.map((row, ri) => (
            <React.Fragment key={ri}>
              <tr className="hover:bg-gray-50">
                <td className="px-4 py-3 border-r border-gray-300 align-top border-b border-gray-300">
                  <input value={row.name} onChange={(e) => changeExerciseName(ri, e.target.value)} placeholder="種目名（例：ベンチプレス）" className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-booking-500" />
                </td>

                {row.sets.slice(0, 5).map((s, si) => (
                  <React.Fragment key={si}>
                    <td className="px-2 py-2 border-r border-gray-200 border-b border-gray-200">
                      <input type="number" value={s.weightKg as any} onChange={(e) => updateCell(ri, si, "weightKg", e.target.value)} className="w-full px-2 py-1 border border-gray-300 rounded text-center focus:outline-none focus:ring-1 focus:ring-booking-500" />
                    </td>
                    <td className="px-2 py-2 border-r border-gray-300 border-b border-gray-200">
                      <input type="number" value={s.reps as any} onChange={(e) => updateCell(ri, si, "reps", e.target.value)} className="w-full px-2 py-1 border border-gray-300 rounded text-center focus:outline-none focus:ring-1 focus:ring-booking-500" />
                    </td>
                  </React.Fragment>
                ))}
              </tr>

              <tr className="border-b border-gray-300 hover:bg-gray-50">
                <td colSpan={11} className="px-4 py-2">
                  <div className="grid grid-cols-5 gap-2">
                    {row.sets.slice(0, 5).map((s, si) => (
                      <input key={si} type="text" value={s.note ?? ""} placeholder={`${si + 1}セット目メモ（任意）`} onChange={(e) => updateSetNote(ri, si, e.target.value)} className="w-full px-2 py-1 text-xs text-gray-600 border border-gray-200 rounded focus:outline-none" />
                    ))}
                  </div>
                </td>
              </tr>
            </React.Fragment>
          ))}

          <tr className="h-24 border-b border-gray-200">
            <td colSpan={12} className="px-4 py-3 text-center text-gray-400">
              <button className="text-booking-600 hover:text-booking-700 font-medium" onClick={addExerciseRow}>
                + 種目を追加
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  );
};

export default TrainingExercisesEditor;
