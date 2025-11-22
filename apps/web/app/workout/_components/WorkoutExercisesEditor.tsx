"use client";
import React from "react";
import type { ExerciseRow } from "../_lib/utils";
import { updateExerciseCell, updateExerciseNote, updateExerciseName, createEmptyExerciseRow } from "../_lib/utils";

type Props = {
  rows: ExerciseRow[];
  onChangeRows: (rows: ExerciseRow[]) => void;
};

const WorkoutExercisesEditor: React.FC<Props> = ({ rows, onChangeRows }) => {
  // rows が undefined の場合は空配列を使用
  const safeRows = rows || [];

  const handleUpdateCell = (ri: number, si: number, key: "weight_kg" | "reps", val: string) => {
    onChangeRows(updateExerciseCell(safeRows, ri, si, key, val));
  };

  const handleUpdateNote = (ri: number, note: string) => {
    onChangeRows(updateExerciseNote(safeRows, ri, note));
  };

  const handleChangeName = (ri: number, name: string) => {
    onChangeRows(updateExerciseName(safeRows, ri, name));
  };

  const handleAddRow = () => {
    onChangeRows([...safeRows, createEmptyExerciseRow()]);
  };

  return (
    <div className="overflow-x-auto">
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
          {/* 最初の3行は常に表示 */}
          {rows?.slice(0, 3).map((row, ri) => (
            <React.Fragment key={ri}>
              <tr className="hover:bg-gray-50">
                <td className="px-4 py-3 border-r border-gray-300 align-top border-b">
                  <input value={row.name} onChange={(e) => handleChangeName(ri, e.target.value)} placeholder="種目名（例：ベンチプレス）" className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-booking-500" />
                </td>

                {row.sets.slice(0, 5).map((s, si) => (
                  <React.Fragment key={si}>
                    <td className="px-2 py-2 border-r border-gray-200 border-b">
                      <div className="flex items-center gap-1">
                        <input type="number" value={s.weight_kg as any} onChange={(e) => handleUpdateCell(ri, si, "weight_kg", e.target.value)} className="w-14 px-2 py-1 border border-gray-300 rounded text-center focus:outline-none focus:ring-1 focus:ring-booking-500" />
                        <span className="text-xs text-gray-600">kg</span>
                      </div>
                    </td>
                    <td className="px-2 py-2 border-r border-gray-300 border-b">
                      <div className="flex items-center gap-1">
                        <input type="number" value={s.reps as any} onChange={(e) => handleUpdateCell(ri, si, "reps", e.target.value)} className="w-14 px-2 py-1 border border-gray-300 rounded text-center focus:outline-none focus:ring-1 focus:ring-booking-500" />
                        <span className="text-xs text-gray-600">rep</span>
                      </div>
                    </td>
                  </React.Fragment>
                ))}
              </tr>

              <tr className="border-b border-gray-300 hover:bg-gray-50">
                <td colSpan={11} className="px-4 py-2">
                  <div className="flex items-center gap-2">
                    <input type="text" value={row.sets[0]?.note ?? ""} placeholder="メモ" onChange={(e) => handleUpdateNote(ri, e.target.value)} className="flex-1 px-3 py-2 text-sm text-gray-700 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-booking-500" />
                  </div>
                </td>
              </tr>
            </React.Fragment>
          ))}

          {/* 4行目以降を表示 */}
          {safeRows.slice(3).map((row, idx) => {
            const ri = idx + 3; // 実際のインデックスは+3
            return (
              <React.Fragment key={ri}>
                <tr className="hover:bg-gray-50">
                  <td className="px-4 py-3 border-r border-gray-300 align-top border-b">
                    <input value={row.name} onChange={(e) => handleChangeName(ri, e.target.value)} placeholder="種目名（例：ベンチプレス）" className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-booking-500" />
                  </td>

                  {row.sets.slice(0, 5).map((s, si) => (
                    <React.Fragment key={si}>
                      <td className="px-2 py-2 border-r border-gray-200 border-b">
                        <div className="flex items-center gap-1">
                          <input type="number" value={s.weight_kg as any} onChange={(e) => handleUpdateCell(ri, si, "weight_kg", e.target.value)} className="w-14 px-2 py-1 border border-gray-300 rounded text-center focus:outline-none focus:ring-1 focus:ring-booking-500" />
                          <span className="text-xs text-gray-600">kg</span>
                        </div>
                      </td>
                      <td className="px-2 py-2 border-r border-gray-300 border-b">
                        <div className="flex items-center gap-1">
                          <input type="number" value={s.reps as any} onChange={(e) => handleUpdateCell(ri, si, "reps", e.target.value)} className="w-14 px-2 py-1 border border-gray-300 rounded text-center focus:outline-none focus:ring-1 focus:ring-booking-500" />
                          <span className="text-xs text-gray-600">rep</span>
                        </div>
                      </td>
                    </React.Fragment>
                  ))}
                </tr>

                <tr className="border-b border-gray-300 hover:bg-gray-50">
                  <td colSpan={11} className="px-4 py-2">
                    <div className="flex items-center gap-2">
                      <input type="text" value={row.sets[0]?.note ?? ""} placeholder="メモ" onChange={(e) => handleUpdateNote(ri, e.target.value)} className="flex-1 px-3 py-2 text-sm text-gray-700 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-booking-500" />
                    </div>
                  </td>
                </tr>
              </React.Fragment>
            );
          })}

          <tr className="h-24 border-b border-gray-200">
            <td colSpan={12} className="px-4 py-3 text-center text-gray-400">
              <button className="text-booking-600 hover:text-booking-700 font-medium" onClick={handleAddRow}>
                + 種目を追加
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  );
};

export default WorkoutExercisesEditor;
