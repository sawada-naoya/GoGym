"use client";
import { useState, useEffect, Fragment } from "react";
import type { ExerciseRow } from "../_lib/utils";
import type { WorkoutPartDTO, ExerciseDTO } from "../_lib/types";
import { updateExerciseCell, updateExerciseNote, createEmptyExerciseRow } from "../_lib/utils";
import ExerciseManageModal from "./ExerciseManageModal";
import { useIsMobile, useMobileExerciseAdjustment } from "../_hooks";

type Props = {
  allExercises: ExerciseRow[]; // 全部位の種目
  onChangeExercises: (exercises: ExerciseRow[]) => void;
  workoutParts: WorkoutPartDTO[];
  selectedPartId: number | null;
  onPartChange: (partId: number | null) => void;
  isUpdate: boolean;
  onSubmit: () => void;
  onRefetchParts: () => void;
  dataKey?: string; // データを識別するキー（日付など）
  onFetchLastRecord: (token: string, exerciseID: number) => Promise<ExerciseDTO | null>;
  token: string;
};

const WorkoutExercisesEditor: React.FC<Props> = ({ allExercises, onChangeExercises, workoutParts, selectedPartId, onPartChange, isUpdate, onSubmit, onRefetchParts, dataKey, onFetchLastRecord, token }) => {
  // 選択中の部位でフィルタリング
  const displayedExercises = selectedPartId ? allExercises.filter((ex) => ex.workout_part_id === selectedPartId) : [];

  // 表示用（最低1つの空種目を保証）
  const exercises = displayedExercises.length > 0 ? displayedExercises : [createEmptyExerciseRow(1)];
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [partExercises, setPartExercises] = useState<Array<{ id: number; name: string; workout_part_id: number | null }>>([]);
  const [lastRecords, setLastRecords] = useState<Map<number, ExerciseDTO>>(new Map());

  // カスタムフック: モバイル判定
  const isMobile = useIsMobile();

  // カスタムフック: モバイル時のセット数調整
  useMobileExerciseAdjustment(isMobile, allExercises, onChangeExercises, dataKey);

  // 部位が変更されたら種目リストを更新
  useEffect(() => {
    if (!selectedPartId) {
      setPartExercises([]);
      return;
    }
    const part = workoutParts.find((workoutPart) => workoutPart.id === selectedPartId);
    setPartExercises(part?.exercises || []);
  }, [selectedPartId, workoutParts]);

  // 種目が選択されたら前回の記録を取得
  useEffect(() => {
    const fetchPreviousRecords = async () => {
      const newLastRecords = new Map<number, ExerciseDTO>();

      for (const exercise of exercises) {
        // 種目にIDがある場合のみ前回記録を取得
        if (exercise.id) {
          try {
            const lastRecord = await onFetchLastRecord(token, exercise.id);
            if (lastRecord) {
              newLastRecords.set(exercise.id, lastRecord);
            }
          } catch (err) {
            console.error(`Failed to fetch last record for exercise ${exercise.id}:`, err);
          }
        }
      }

      setLastRecords(newLastRecords);
    };

    fetchPreviousRecords();
  }, [exercises.map((e) => e.id).join(","), token, onFetchLastRecord]);

  // 表示されているexercise（フィルタ済み）を更新して、全体を再構築
  const updateDisplayedExercise = (updatedDisplayed: ExerciseRow[]) => {
    const otherExercises = allExercises.filter((ex) => ex.workout_part_id !== selectedPartId);
    onChangeExercises([...otherExercises, ...updatedDisplayed]);
  };

  const handleUpdateCell = (exerciseIndex: number, setIndex: number, key: "weight_kg" | "reps", value: string) => {
    const updated = updateExerciseCell(exercises, exerciseIndex, setIndex, key, value);
    updateDisplayedExercise(updated);
  };

  const handleUpdateNote = (exerciseIndex: number, note: string) => {
    const updated = updateExerciseNote(exercises, exerciseIndex, note);
    updateDisplayedExercise(updated);
  };

  const handleRemoveSet = (exerciseIndex: number, setIndex: number) => {
    if (exercises[exerciseIndex].sets.length === 1) return;
    const updatedExercises = structuredClone(exercises);
    updatedExercises[exerciseIndex].sets.splice(setIndex, 1);
    updateDisplayedExercise(updatedExercises);
  };

  const handleChangeExerciseName = (exerciseIndex: number, exerciseName: string) => {
    const selectedExercise = partExercises.find((exercise) => exercise.name === exerciseName);
    const updatedExercises = structuredClone(exercises);
    updatedExercises[exerciseIndex].name = exerciseName;
    updatedExercises[exerciseIndex].id = selectedExercise?.id || null;
    updatedExercises[exerciseIndex].workout_part_id = selectedExercise?.workout_part_id || selectedPartId || null;
    updateDisplayedExercise(updatedExercises);
  };

  const handleAddExerciseRow = () => {
    const newExercise = createEmptyExerciseRow(1);
    newExercise.workout_part_id = selectedPartId;
    updateDisplayedExercise([...exercises, newExercise]);
  };

  const handleAddSet = (exerciseIndex: number) => {
    const updatedExercises = structuredClone(exercises);
    const newSetNumber = updatedExercises[exerciseIndex].sets.length + 1;
    updatedExercises[exerciseIndex].sets.push({
      set_number: newSetNumber,
      weight_kg: 0,
      reps: 0,
      note: null,
    });
    updateDisplayedExercise(updatedExercises);
  };

  const handleRemoveExercise = (exerciseIndex: number) => {
    if (exercises.length === 1) return;
    updateDisplayedExercise(exercises.filter((_, index) => index !== exerciseIndex));
  };

  const handleCopyLastRecord = (exerciseIndex: number) => {
    const currentExercise = exercises[exerciseIndex];
    if (!currentExercise.id) return;

    const previousRecord = lastRecords.get(currentExercise.id);
    if (!previousRecord || !previousRecord.sets || previousRecord.sets.length === 0) return;

    const updatedExercises = structuredClone(exercises);

    // 前回記録のセット数に合わせてセットを調整
    const previousSets = previousRecord.sets;
    updatedExercises[exerciseIndex].sets = previousSets.map((previousSet, index) => ({
      id: null, // 新規セットなのでIDはnull
      set_number: index + 1,
      weight_kg: String(previousSet.weight_kg || ""),
      reps: String(previousSet.reps || ""),
      note: previousSet.note || null,
    }));

    updateDisplayedExercise(updatedExercises);
  };

  const handleCopySetBelow = (exerciseIndex: number, setIndex: number) => {
    const currentExercise = exercises[exerciseIndex];
    if (!currentExercise.sets || currentExercise.sets.length >= 5) return; // 最大5セット

    const currentSet = currentExercise.sets[setIndex];

    const updatedExercises = structuredClone(exercises);

    // 現在のセットの内容をコピーして直後に追加
    const copiedSet = {
      id: null, // 新規セットなのでIDはnull
      set_number: setIndex + 2, // 一時的なセット番号（後で再計算される）
      weight_kg: currentSet.weight_kg || 0,
      reps: currentSet.reps || 0,
      note: null, // メモはコピーしない
    };

    // 押されたセットの直後（setIndex + 1の位置）に挿入
    updatedExercises[exerciseIndex].sets.splice(setIndex + 1, 0, copiedSet);

    // セット番号を再計算（1から順番に振り直す）
    updatedExercises[exerciseIndex].sets = updatedExercises[exerciseIndex].sets.map((set, index) => ({
      ...set,
      set_number: index + 1,
    }));

    updateDisplayedExercise(updatedExercises);
  };

  const handlePartChange = (idStr: string) => {
    if (!idStr) {
      onPartChange(null);
      return;
    }
    const numId = Number(idStr);
    onPartChange(numId);
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
            <select value={selectedPartId?.toString() ?? ""} onChange={(e) => handlePartChange(e.target.value)} className="flex-1 px-2 py-1.5 text-xs border border-gray-300 rounded-lg focus:outline-none focus:ring-1 focus:ring-booking-500 focus:border-booking-500 bg-white transition-colors">
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
        {exercises.map((exercise, exerciseIndex) => (
          <div key={exerciseIndex} className="bg-white rounded-xl shadow-sm border border-gray-100 p-3">
            {/* 種目名とボタン */}
            <div className="flex items-center gap-2 mb-2">
              <select value={exercise.name} onChange={(e) => handleChangeExerciseName(exerciseIndex, e.target.value)} className="flex-1 px-2 py-1.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-1 focus:ring-booking-500 focus:border-booking-500 bg-white text-xs font-semibold transition-colors">
                <option value="">種目を選択</option>
                {partExercises.map((partExercise) => (
                  <option key={partExercise.id} value={partExercise.name}>
                    {partExercise.name}
                  </option>
                ))}
              </select>
              {exercise.id && lastRecords.get(exercise.id) && (
                <button type="button" onClick={() => handleCopyLastRecord(exerciseIndex)} className="p-1.5 text-booking-600 hover:text-booking-700 bg-booking-50 hover:bg-booking-100 rounded-lg transition-colors border border-booking-200" title="前回記録をコピー">
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                  </svg>
                </button>
              )}
              <button type="button" onClick={() => handleRemoveExercise(exerciseIndex)} className="p-1 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors" disabled={exercises.length === 1}>
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>

            {/* セットリスト（縦並び） */}
            <div className="space-y-1 mb-1.5">
              {exercise.sets.map((set, setIndex) => {
                const previousRecord = exercise.id ? lastRecords.get(exercise.id) : null;
                const previousSet = previousRecord?.sets?.[setIndex];

                return (
                  <div key={setIndex} className="flex items-center gap-1 bg-gradient-to-r from-gray-50 to-gray-50/50 rounded-lg p-1.5 border border-gray-100">
                    <span className="text-xs font-bold text-gray-600 w-4 flex-shrink-0">{setIndex + 1}</span>
                    <input
                      type="number"
                      value={set.weight_kg as any}
                      onChange={(e) => handleUpdateCell(exerciseIndex, setIndex, "weight_kg", e.target.value)}
                      className="w-12 px-1 py-1 text-xs font-semibold border border-gray-300 rounded-md text-center focus:outline-none focus:ring-1 focus:ring-booking-500 focus:border-booking-500 bg-white transition-colors [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
                    />
                    <span className="text-[9px] font-medium text-gray-600">kg</span>
                    <span className="text-gray-300 text-xs font-bold">×</span>
                    <input
                      type="number"
                      value={set.reps as any}
                      onChange={(e) => handleUpdateCell(exerciseIndex, setIndex, "reps", e.target.value)}
                      className="w-12 px-1 py-1 text-xs font-semibold border border-gray-300 rounded-md text-center focus:outline-none focus:ring-1 focus:ring-booking-500 focus:border-booking-500 bg-white transition-colors [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
                    />
                    <span className="text-[9px] font-medium text-gray-600">回</span>
                    {previousSet && (
                      <span className="text-[9px] text-gray-500 ml-1">
                        前回 : {previousSet.weight_kg}kg×{previousSet.reps}
                      </span>
                    )}
                    <div className="ml-auto flex items-center gap-0.5">
                      <button type="button" onClick={() => handleCopySetBelow(exerciseIndex, setIndex)} className="p-0.5 text-purple-600 hover:text-purple-700 bg-purple-50 hover:bg-purple-100 rounded transition-colors border border-purple-200 disabled:opacity-30 disabled:cursor-not-allowed" disabled={exercise.sets.length >= 5} title="下にコピー">
                        <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2.5}>
                          <path strokeLinecap="round" strokeLinejoin="round" d="M15 13l-3 3m0 0l-3-3m3 3V8m0 13a9 9 0 110-18 9 9 0 010 18z" />
                        </svg>
                      </button>
                      <button type="button" onClick={() => handleRemoveSet(exerciseIndex, setIndex)} className="p-0.5 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded transition-colors disabled:opacity-30 disabled:cursor-not-allowed" disabled={exercise.sets.length === 1}>
                        <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2.5}>
                          <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                      </button>
                    </div>
                  </div>
                );
              })}
            </div>

            {/* セット追加ボタン（5セット未満の場合のみ表示） */}
            {exercise.sets.length < 5 && (
              <div className="flex justify-center mb-1.5">
                <button type="button" onClick={() => handleAddSet(exerciseIndex)} className="p-1 text-booking-600 hover:text-booking-700 bg-booking-50 hover:bg-booking-100 rounded-lg border border-booking-200 transition-colors">
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
                value={exercise.sets[0]?.note ?? ""}
                placeholder="メモ"
                onChange={(e) => handleUpdateNote(exerciseIndex, e.target.value)}
                className="w-full px-2 py-1.5 text-xs text-gray-700 border border-gray-300 rounded-lg focus:outline-none focus:ring-1 focus:ring-booking-500 focus:border-booking-500 bg-gray-50 transition-colors placeholder:text-gray-400"
              />
            </div>
          </div>
        ))}

        {/* 種目追加ボタン */}
        <div className="flex justify-center pt-0.5">
          <button type="button" onClick={handleAddExerciseRow} className="flex items-center gap-1.5 px-4 py-2 text-xs text-booking-600 font-bold hover:text-booking-700 bg-booking-50 hover:bg-booking-100 rounded-xl border border-booking-200 hover:border-booking-300 transition-all shadow-sm hover:shadow">
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
            <select value={selectedPartId?.toString() ?? ""} onChange={(e) => handlePartChange(e.target.value)} className="w-40 px-3 py-2.5 text-base border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-booking-500 bg-white">
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
                {[1, 2, 3, 4, 5].map((setNumber) => (
                  <th key={setNumber} colSpan={2} className="px-4 py-3 text-center font-medium text-gray-700 border-r border-gray-300">
                    {setNumber}セット
                  </th>
                ))}
                <th className="px-4 py-3 text-center font-medium text-gray-700 w-12"></th>
              </tr>
              <tr className="border-b border-gray-300 bg-gray-50">
                <th className="px-4 py-2 border-r border-gray-300" />
                {[1, 2, 3, 4, 5].map((setNumber) => (
                  <Fragment key={setNumber}>
                    <th className="px-2 py-2 text-center text-sm font-medium text-gray-600 border-r border-gray-200">重量</th>
                    <th className="px-2 py-2 text-center text-sm font-medium text-gray-600 border-r border-gray-300">回数</th>
                  </Fragment>
                ))}
                <th className="px-2 py-2 border-r border-gray-300"></th>
              </tr>
            </thead>

            <tbody>
              {exercises.map((exercise, exerciseIndex) => (
                <Fragment key={exerciseIndex}>
                  <tr className="hover:bg-gray-50">
                    <td className="px-4 py-3 border-r border-gray-300 align-top border-b">
                      <div className="flex items-center gap-2">
                        <select value={exercise.name} onChange={(e) => handleChangeExerciseName(exerciseIndex, e.target.value)} className="flex-1 px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-booking-500 bg-white text-sm truncate">
                          <option value="">種目を選択</option>
                          {partExercises.map((partExercise) => (
                            <option key={partExercise.id} value={partExercise.name}>
                              {partExercise.name}
                            </option>
                          ))}
                        </select>
                        {exercise.id && lastRecords.get(exercise.id) && (
                          <button type="button" onClick={() => handleCopyLastRecord(exerciseIndex)} className="p-1 text-booking-600 hover:text-booking-700 hover:bg-booking-50 rounded transition-colors flex-shrink-0" title="前回記録をコピー">
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                            </svg>
                          </button>
                        )}
                      </div>
                    </td>

                    {[0, 1, 2, 3, 4].map((setIndex) => {
                      const set = exercise.sets[setIndex];
                      const previousRecord = exercise.id ? lastRecords.get(exercise.id) : null;
                      const previousSet = previousRecord?.sets?.[setIndex];

                      return (
                        <Fragment key={setIndex}>
                          <td className="px-2 py-2 border-r border-gray-200 border-b">
                            {set ? (
                              <div className="flex flex-col gap-0.5">
                                <div className="flex items-center gap-1">
                                  <input
                                    type="number"
                                    value={set.weight_kg as any}
                                    onChange={(e) => handleUpdateCell(exerciseIndex, setIndex, "weight_kg", e.target.value)}
                                    className="w-14 px-2 py-1 border border-gray-300 rounded text-center focus:outline-none focus:ring-1 focus:ring-booking-500 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
                                  />
                                  <span className="text-xs text-gray-600">kg</span>
                                </div>
                                {previousSet && <span className="text-[10px] text-gray-500 text-center">前回: {previousSet.weight_kg}kg</span>}
                              </div>
                            ) : (
                              <div className="h-8"></div>
                            )}
                          </td>
                          <td className="px-2 py-2 border-r border-gray-300 border-b">
                            {set ? (
                              <div className="flex flex-col gap-0.5">
                                <div className="flex items-center gap-1">
                                  <input
                                    type="number"
                                    value={set.reps as any}
                                    onChange={(e) => handleUpdateCell(exerciseIndex, setIndex, "reps", e.target.value)}
                                    className="w-14 px-2 py-1 border border-gray-300 rounded text-center focus:outline-none focus:ring-1 focus:ring-booking-500 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
                                  />
                                  <span className="text-xs text-gray-600">rep</span>
                                </div>
                                {previousSet && <span className="text-[10px] text-gray-500 text-center">前回: {previousSet.reps}回</span>}
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
                      {exercise.sets.length < 5 && (
                        <button type="button" onClick={() => handleAddSet(exerciseIndex)} className="p-1 text-booking-600 hover:text-booking-700 hover:bg-booking-50 rounded transition-colors">
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
                        <input type="text" value={exercise.sets[0]?.note ?? ""} placeholder="メモ" onChange={(e) => handleUpdateNote(exerciseIndex, e.target.value)} className="flex-1 px-3 py-2 text-sm text-gray-700 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-booking-500 truncate" />
                      </div>
                    </td>
                  </tr>
                </Fragment>
              ))}

              <tr className="h-24 border-b border-gray-200">
                <td colSpan={12} className="px-4 py-3 text-center text-gray-400">
                  <button className="inline-flex items-center justify-center w-10 h-10 text-booking-600 hover:text-booking-700 hover:bg-booking-50 rounded-full transition-colors" onClick={handleAddExerciseRow}>
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
