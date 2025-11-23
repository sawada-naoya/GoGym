"use client";
import React, { useState } from "react";
import type { ExerciseRow } from "../_lib/utils";
import type { WorkoutPartDTO, WorkoutFormDTO } from "../_lib/types";
import { updateExerciseCell, updateExerciseNote, updateExerciseName, createEmptyExerciseRow } from "../_lib/utils";
import ExerciseManageModal from "./ExerciseManageModal";

type Props = {
  workoutExercises: ExerciseRow[];
  onChangeExercises: (exercises: ExerciseRow[]) => void;
  workoutParts: WorkoutPartDTO[];
  selectedPart: WorkoutFormDTO["workout_part"];
  onPartChange: (part: WorkoutFormDTO["workout_part"]) => void;
  isUpdate: boolean;
  onSubmit: () => void;
};

const WorkoutExercisesEditor: React.FC<Props> = ({ workoutExercises, onChangeExercises, workoutParts, selectedPart, onPartChange, isUpdate, onSubmit }) => {
  // workoutExercises が undefined の場合は空配列を使用
  const safeExercises = workoutExercises || [];
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [partExercises, setPartExercises] = useState<Array<{ id: number; name: string }>>([]);

  // 部位が変更されたら種目リストを更新
  React.useEffect(() => {
    if (!selectedPart?.id) {
      setPartExercises([]);
      return;
    }
    const part = workoutParts.find((p) => p.id === selectedPart.id);
    setPartExercises(part?.exercises || []);
  }, [selectedPart, workoutParts]);

  const handleUpdateCell = (ri: number, si: number, key: "weight_kg" | "reps", val: string) => {
    onChangeExercises(updateExerciseCell(safeExercises, ri, si, key, val));
  };

  const handleUpdateNote = (ri: number, note: string) => {
    onChangeExercises(updateExerciseNote(safeExercises, ri, note));
  };

  const handleChangeName = (ri: number, name: string) => {
    onChangeExercises(updateExerciseName(safeExercises, ri, name));
  };

  const handleAddRow = () => {
    onChangeExercises([...safeExercises, createEmptyExerciseRow()]);
  };

  const handlePartChange = (idStr: string) => {
    if (!idStr) {
      onPartChange({ id: null, name: null, source: null });
      return;
    }
    const numId = Number(idStr);
    const part = workoutParts.find((p) => p.id === numId);
    if (part) {
      onPartChange({
        id: part.id,
        name: part.name,
        source: "custom", // すべてカスタム扱い（is_default削除により）
      });
    }
  };

  // 種目登録成功時にページをリロード
  const handleSuccess = () => {
    window.location.reload();
  };

  return (
    <div className="bg-white rounded-lg shadow mb-6 p-6">
      {/* ヘッダー: 部位選択とアクションボタン */}
      <div className="flex items-center justify-between gap-3 mb-6">
        <div className="flex items-center gap-2">
          <label className="text-sm font-medium text-gray-700">部位</label>
          <select value={selectedPart?.id?.toString() ?? ""} onChange={(e) => handlePartChange(e.target.value)} className="w-40 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500 bg-white">
            <option value=""></option>
            {workoutParts.map((part) => (
              <option key={part.id} value={part.id.toString()}>
                {part.name}
              </option>
            ))}
          </select>
        </div>
        <div className="flex gap-2">
          <button type="button" onClick={() => setIsModalOpen(true)} className="px-4 py-2 rounded-md bg-booking-50 text-booking-700 hover:bg-booking-100 transition-colors whitespace-nowrap text-sm border border-booking-300">
            + 種目追加
          </button>
          <button onClick={onSubmit} className="px-5 py-2 rounded-md bg-booking-600 text-white hover:bg-booking-700 transition-colors disabled:opacity-50 whitespace-nowrap">
            {isUpdate ? "修正" : "登録"}
          </button>
        </div>
      </div>

      {/* 種目追加モーダル */}
      <ExerciseManageModal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} workoutParts={workoutParts} onSuccess={handleSuccess} />

      {/* テーブル */}
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
            {workoutExercises?.map((row, ri) => (
              <React.Fragment key={ri}>
                <tr className="hover:bg-gray-50">
                  <td className="px-4 py-3 border-r border-gray-300 align-top border-b">
                    <select value={row.name} onChange={(e) => handleChangeName(ri, e.target.value)} className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-booking-500 bg-white">
                      <option value="">種目を選択</option>
                      {partExercises.map((ex) => (
                        <option key={ex.id} value={ex.name}>
                          {ex.name}
                        </option>
                      ))}
                    </select>
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
            {safeExercises.slice(3).map((row, idx) => {
              const ri = idx + 3; // 実際のインデックスは+3
              return (
                <React.Fragment key={ri}>
                  <tr className="hover:bg-gray-50">
                    <td className="px-4 py-3 border-r border-gray-300 align-top border-b">
                      <select value={row.name} onChange={(e) => handleChangeName(ri, e.target.value)} className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-booking-500 bg-white">
                        <option value="">種目を選択</option>
                        {partExercises.map((ex) => (
                          <option key={ex.id} value={ex.name}>
                            {ex.name}
                          </option>
                        ))}
                      </select>
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
                <button className="inline-flex items-center justify-center w-10 h-10 text-booking-600 hover:text-booking-700 hover:bg-booking-50 rounded-full transition-colors" onClick={handleAddRow}>
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                  </svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default WorkoutExercisesEditor;
