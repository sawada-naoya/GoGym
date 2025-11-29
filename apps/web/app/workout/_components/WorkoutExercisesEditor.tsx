"use client";
import { useState, useEffect, Fragment } from "react";
import type { ExerciseRow } from "../_lib/utils";
import type { WorkoutPartDTO, WorkoutFormDTO } from "../_lib/types";
import { updateExerciseCell, updateExerciseNote, createEmptyExerciseRow } from "../_lib/utils";
import ExerciseManageModal from "./ExerciseManageModal";
import { useIsMobile, useMobileExerciseAdjustment } from "../_hooks";

type Props = {
  workoutExercises: ExerciseRow[];
  onChangeExercises: (exercises: ExerciseRow[]) => void;
  workoutParts: WorkoutPartDTO[];
  selectedPart: WorkoutFormDTO["workout_part"];
  onPartChange: (part: WorkoutFormDTO["workout_part"]) => void;
  isUpdate: boolean;
  onSubmit: () => void;
  onRefetchParts: () => void;
  dataKey?: string; // データを識別するキー（日付など）
};

const WorkoutExercisesEditor: React.FC<Props> = ({ workoutExercises, onChangeExercises, workoutParts, selectedPart, onPartChange, isUpdate, onSubmit, onRefetchParts, dataKey }) => {
  // workoutExercises が undefined または空の場合は1つの空種目を用意
  const safeExercises = workoutExercises && workoutExercises.length > 0 ? workoutExercises : [createEmptyExerciseRow(1)];
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [partExercises, setPartExercises] = useState<Array<{ id: number; name: string; workout_part_id: number | null }>>([]);

  // カスタムフック: モバイル判定
  const isMobile = useIsMobile();

  // カスタムフック: モバイル時のセット数調整
  useMobileExerciseAdjustment(isMobile, workoutExercises, onChangeExercises, dataKey);

  // 部位が変更されたら種目リストを更新
  useEffect(() => {
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

  const handleRemoveSet = (ri: number, si: number) => {
    if (safeExercises[ri].sets.length === 1) return;
    const next = structuredClone(safeExercises);
    next[ri].sets.splice(si, 1);
    onChangeExercises(next);
  };

  const handleChangeName = (ri: number, name: string) => {
    // 選択された種目のIDと部位IDを取得
    const selectedExercise = partExercises.find((ex) => ex.name === name);
    const next = structuredClone(safeExercises);
    next[ri].name = name;
    next[ri].id = selectedExercise?.id || null;
    next[ri].workout_part_id = selectedExercise?.workout_part_id || selectedPart?.id || null;
    onChangeExercises(next);
  };

  const handleAddRow = () => {
    // 新規種目は1セットから開始
    onChangeExercises([...safeExercises, createEmptyExerciseRow(1)]);
  };

  const handleAddSet = (ri: number) => {
    const next = structuredClone(safeExercises);
    const newSetNumber = next[ri].sets.length + 1;
    next[ri].sets.push({
      set_number: newSetNumber,
      weight_kg: "",
      reps: "",
      note: null,
    });
    onChangeExercises(next);
  };

  const handleRemoveExercise = (ri: number) => {
    if (safeExercises.length === 1) return;
    onChangeExercises(safeExercises.filter((_, i) => i !== ri));
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

  // 種目登録成功時に部位データを再取得
  const handleSuccess = () => {
    onRefetchParts();
  };

  return (
    <>
      {/* 種目追加モーダル */}
      <ExerciseManageModal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} workoutParts={workoutParts} onSuccess={handleSuccess} />

      {/* モバイル: 部位選択エリア */}
      <div className="md:hidden bg-white rounded-xl shadow-sm border border-gray-100 mb-2 p-3">
        <div className="flex items-center gap-2">
          <div className="flex items-center gap-2 flex-1">
            <label className="text-xs font-semibold text-gray-700 whitespace-nowrap">部位</label>
            <select value={selectedPart?.id?.toString() ?? ""} onChange={(e) => handlePartChange(e.target.value)} className="flex-1 px-2 py-1.5 text-xs border-2 border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-booking-500 focus:border-booking-500 bg-white transition-colors">
              <option value="">選択してください</option>
              {workoutParts.map((part) => (
                <option key={part.id} value={part.id.toString()}>
                  {part.name}
                </option>
              ))}
            </select>
          </div>
          <button type="button" onClick={() => setIsModalOpen(true)} className="px-2.5 py-1.5 rounded-lg text-booking-600 bg-booking-50 hover:bg-booking-100 transition-colors text-xs border border-booking-200 font-semibold whitespace-nowrap">
            種目作成
          </button>
        </div>
      </div>

      {/* モバイル: 種目入力エリア（各種目ごとにカード） */}
      <div className="md:hidden space-y-2">
        {safeExercises.map((row, ri) => (
          <div key={ri} className="bg-white rounded-xl shadow-sm border border-gray-100 p-3">
            {/* 種目名と削除ボタン */}
            <div className="flex items-center justify-between mb-2">
              <select value={row.name} onChange={(e) => handleChangeName(ri, e.target.value)} className="flex-1 px-2 py-1.5 border-2 border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-booking-500 focus:border-booking-500 bg-white text-xs font-semibold transition-colors">
                <option value="">種目を選択</option>
                {partExercises.map((ex) => (
                  <option key={ex.id} value={ex.name}>
                    {ex.name}
                  </option>
                ))}
              </select>
              <button type="button" onClick={() => handleRemoveExercise(ri)} className="ml-2 p-1 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors" disabled={safeExercises.length === 1}>
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>

            {/* セットリスト（縦並び） */}
            <div className="space-y-1 mb-1.5">
              {row.sets.map((s, si) => (
                <div key={si} className="flex items-center gap-1.5 bg-gradient-to-r from-gray-50 to-gray-50/50 rounded-lg p-1.5 border border-gray-100">
                  <span className="text-xs font-bold text-gray-600 w-5">{si + 1}</span>
                  <input
                    type="number"
                    value={s.weight_kg as any}
                    onChange={(e) => handleUpdateCell(ri, si, "weight_kg", e.target.value)}
                    placeholder="0"
                    className="w-14 px-1.5 py-1 text-xs font-semibold border-2 border-gray-200 rounded-md text-center focus:outline-none focus:ring-2 focus:ring-booking-400 focus:border-booking-400 bg-white transition-colors [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
                  />
                  <span className="text-[10px] font-medium text-gray-600">kg</span>
                  <span className="text-gray-300 text-xs font-bold">×</span>
                  <input
                    type="number"
                    value={s.reps as any}
                    onChange={(e) => handleUpdateCell(ri, si, "reps", e.target.value)}
                    placeholder="0"
                    className="w-14 px-1.5 py-1 text-xs font-semibold border-2 border-gray-200 rounded-md text-center focus:outline-none focus:ring-2 focus:ring-booking-400 focus:border-booking-400 bg-white transition-colors [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
                  />
                  <span className="text-[10px] font-medium text-gray-600">回</span>
                  <button type="button" onClick={() => handleRemoveSet(ri, si)} className="ml-auto p-1 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded-md transition-colors" disabled={row.sets.length === 1}>
                    <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2.5}>
                      <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
              ))}
            </div>

            {/* セット追加ボタン（5セット未満の場合のみ表示） */}
            {row.sets.length < 5 && (
              <div className="flex justify-center mb-1.5">
                <button type="button" onClick={() => handleAddSet(ri)} className="p-1 text-booking-600 hover:text-booking-700 bg-booking-50 hover:bg-booking-100 rounded-lg border border-booking-200 transition-colors">
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
                  </svg>
                </button>
              </div>
            )}

            {/* メモ */}
            <div>
              <input
                type="text"
                value={row.sets[0]?.note ?? ""}
                placeholder="メモ"
                onChange={(e) => handleUpdateNote(ri, e.target.value)}
                className="w-full px-2 py-1.5 text-xs text-gray-700 border-2 border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-booking-400 focus:border-booking-400 bg-gray-50 transition-colors placeholder:text-gray-400"
              />
            </div>
          </div>
        ))}

        {/* 種目追加ボタン */}
        <div className="flex justify-center pt-0.5">
          <button type="button" onClick={handleAddRow} className="flex items-center gap-1.5 px-4 py-2 text-xs text-booking-600 font-bold hover:text-booking-700 bg-booking-50 hover:bg-booking-100 rounded-xl border-2 border-booking-200 hover:border-booking-300 transition-all shadow-sm hover:shadow">
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
            </svg>
            <span>種目を追加</span>
          </button>
        </div>
      </div>

      {/* デスクトップ: ヘッダーとテーブル */}
      <div className="hidden md:block bg-white rounded-lg shadow mb-4 md:mb-6 p-3 md:p-6">
        {/* ヘッダー: 部位選択と種目作成ボタン、登録/更新ボタン */}
        <div className="flex items-center justify-between gap-2 mb-6">
          <div className="flex items-center gap-2">
            <label className="text-base font-medium text-gray-700">部位</label>
            <select value={selectedPart?.id?.toString() ?? ""} onChange={(e) => handlePartChange(e.target.value)} className="w-40 px-3 py-2.5 text-base border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500 bg-white">
              <option value="">選択してください</option>
              {workoutParts.map((part) => (
                <option key={part.id} value={part.id.toString()}>
                  {part.name}
                </option>
              ))}
            </select>
            <button type="button" onClick={() => setIsModalOpen(true)} className="px-4 py-2.5 rounded-md text-booking-600 hover:bg-booking-50 transition-colors whitespace-nowrap text-base border border-gray-200 hover:border-booking-300 font-medium">
              種目作成
            </button>
          </div>

          <button onClick={onSubmit} className="px-8 py-2.5 text-base rounded-md bg-booking-600 text-white hover:bg-booking-700 transition-colors font-medium">
            {isUpdate ? "更新" : "登録"}
          </button>
        </div>

        {/* テーブル表示 */}
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
                <th className="px-4 py-3 text-center font-medium text-gray-700 w-12"></th>
              </tr>
              <tr className="border-b border-gray-300 bg-gray-50">
                <th className="px-4 py-2 border-r border-gray-300" />
                {[1, 2, 3, 4, 5].map((n) => (
                  <Fragment key={n}>
                    <th className="px-2 py-2 text-center text-sm font-medium text-gray-600 border-r border-gray-200">重量</th>
                    <th className="px-2 py-2 text-center text-sm font-medium text-gray-600 border-r border-gray-300">回数</th>
                  </Fragment>
                ))}
                <th className="px-2 py-2 border-r border-gray-300"></th>
              </tr>
            </thead>

            <tbody>
              {safeExercises.map((row, ri) => (
                <Fragment key={ri}>
                  <tr className="hover:bg-gray-50">
                    <td className="px-4 py-3 border-r border-gray-300 align-top border-b">
                      <select value={row.name} onChange={(e) => handleChangeName(ri, e.target.value)} className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-booking-500 bg-white text-sm truncate">
                        <option value="">種目を選択</option>
                        {partExercises.map((ex) => (
                          <option key={ex.id} value={ex.name}>
                            {ex.name}
                          </option>
                        ))}
                      </select>
                    </td>

                    {[0, 1, 2, 3, 4].map((si) => {
                      const s = row.sets[si];
                      return (
                        <Fragment key={si}>
                          <td className="px-2 py-2 border-r border-gray-200 border-b">
                            {s ? (
                              <div className="flex items-center gap-1">
                                <input
                                  type="number"
                                  value={s.weight_kg as any}
                                  onChange={(e) => handleUpdateCell(ri, si, "weight_kg", e.target.value)}
                                  className="w-14 px-2 py-1 border border-gray-300 rounded text-center focus:outline-none focus:ring-1 focus:ring-booking-500 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
                                />
                                <span className="text-xs text-gray-600">kg</span>
                              </div>
                            ) : (
                              <div className="h-8"></div>
                            )}
                          </td>
                          <td className="px-2 py-2 border-r border-gray-300 border-b">
                            {s ? (
                              <div className="flex items-center gap-1">
                                <input
                                  type="number"
                                  value={s.reps as any}
                                  onChange={(e) => handleUpdateCell(ri, si, "reps", e.target.value)}
                                  className="w-14 px-2 py-1 border border-gray-300 rounded text-center focus:outline-none focus:ring-1 focus:ring-booking-500 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
                                />
                                <span className="text-xs text-gray-600">rep</span>
                              </div>
                            ) : (
                              <div className="h-8"></div>
                            )}
                          </td>
                        </Fragment>
                      );
                    })}

                    {/* セット追加ボタン（5セット未満の場合のみ表示） */}
                    <td className="px-2 py-2 text-center border-r border-gray-300 border-b">
                      {row.sets.length < 5 && (
                        <button type="button" onClick={() => handleAddSet(ri)} className="p-1 text-booking-600 hover:text-booking-700 hover:bg-booking-50 rounded transition-colors">
                          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                          </svg>
                        </button>
                      )}
                    </td>
                  </tr>

                  <tr className="border-b border-gray-300 hover:bg-gray-50">
                    <td colSpan={12} className="px-4 py-2">
                      <div className="flex items-center gap-2">
                        <input type="text" value={row.sets[0]?.note ?? ""} placeholder="メモ" onChange={(e) => handleUpdateNote(ri, e.target.value)} className="flex-1 px-3 py-2 text-sm text-gray-700 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-booking-500 truncate" />
                      </div>
                    </td>
                  </tr>
                </Fragment>
              ))}

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

      {/* 登録/更新ボタン（モバイル: 下部固定） */}
      <div className="md:hidden fixed bottom-0 left-0 right-0 p-3 bg-white border-t border-gray-200 shadow-lg z-10">
        <button onClick={onSubmit} className="w-full py-2.5 text-sm rounded-xl bg-booking-600 text-white hover:bg-booking-700 active:bg-booking-800 transition-colors font-bold shadow-sm">
          {isUpdate ? "更新する" : "登録する"}
        </button>
      </div>
    </>
  );
};

export default WorkoutExercisesEditor;
