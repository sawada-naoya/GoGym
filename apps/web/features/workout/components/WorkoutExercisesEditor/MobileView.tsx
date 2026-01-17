"use client";
import { memo } from "react";
import type { ExerciseRow } from "@/features/workout/lib/utils";
import type { WorkoutPartDTO, ExerciseDTO } from "@/types/workout";
import { getLocalizedPartName } from "@/features/workout/lib/utils";

type MobileViewProps = {
  t: any;
  i18n: any;
  exercises: ExerciseRow[];
  workoutParts: WorkoutPartDTO[];
  selectedPartId: number | null;
  partExercises: Array<{
    id: number;
    name: string;
    workout_part_id: number | null;
  }>;
  lastRecords: Map<number, ExerciseDTO>;
  isUpdate: boolean;
  isModalOpen: boolean;
  setIsModalOpen: (open: boolean) => void;
  handlePartChange: (idStr: string) => void;
  handleChangeExerciseName: (
    exerciseIndex: number,
    exerciseName: string,
  ) => void;
  handleCopyLastRecord: (exerciseIndex: number) => void;
  handleRemoveExercise: (exerciseIndex: number) => void;
  handleUpdateCell: (
    exerciseIndex: number,
    setIndex: number,
    key: "weight_kg" | "reps",
    value: string,
  ) => void;
  handleCopySetBelow: (exerciseIndex: number, setIndex: number) => void;
  handleRemoveSet: (exerciseIndex: number, setIndex: number) => void;
  handleAddSet: (exerciseIndex: number) => void;
  handleUpdateNote: (exerciseIndex: number, note: string) => void;
  handleAddExerciseRow: () => void;
  onSubmit: () => void;
  handleSuccess: () => void;
};

const MobileView = memo<MobileViewProps>(
  ({
    t,
    i18n,
    exercises,
    workoutParts,
    selectedPartId,
    partExercises,
    lastRecords,
    isUpdate,
    setIsModalOpen,
    handlePartChange,
    handleChangeExerciseName,
    handleCopyLastRecord,
    handleRemoveExercise,
    handleUpdateCell,
    handleCopySetBelow,
    handleRemoveSet,
    handleAddSet,
    handleUpdateNote,
    handleAddExerciseRow,
    onSubmit,
  }) => {
    return (
      <>
        {/* モバイル: 部位選択エリア */}
        <div className="md:hidden bg-white rounded-xl shadow-sm border border-gray-100 mb-2 p-3">
          <div className="flex items-center gap-2">
            <div className="flex items-center gap-2 flex-1">
              <label className="text-xs font-semibold text-gray-700 whitespace-nowrap">
                {t("workout.exercises.partLabel")}
              </label>
              <select
                value={selectedPartId?.toString() ?? ""}
                onChange={(e) => handlePartChange(e.target.value)}
                className="flex-1 px-2 py-1.5 text-xs border border-gray-300 rounded-lg focus:outline-none focus:ring-1 focus:ring-booking-500 focus:border-booking-500 bg-white transition-colors"
              >
                <option value="">
                  {t("workout.exercises.partPlaceholder")}
                </option>
                {workoutParts.map((part) => (
                  <option key={part.id} value={part.id.toString()}>
                    {getLocalizedPartName(part, i18n.language)}
                  </option>
                ))}
              </select>
            </div>
            <button
              type="button"
              onClick={() => setIsModalOpen(true)}
              className="px-2.5 py-1.5 rounded-lg text-booking-600 bg-booking-50 hover:bg-booking-100 transition-colors text-xs border border-booking-200 font-semibold whitespace-nowrap"
            >
              {t("workout.exercises.createExerciseButton")}
            </button>
          </div>
        </div>

        {/* モバイル: 種目入力エリア（各種目ごとにカード） */}
        <div className="md:hidden space-y-2">
          {exercises.map((exercise, exerciseIndex) => (
            <div
              key={exerciseIndex}
              className="bg-white rounded-xl shadow-sm border border-gray-100 p-3"
            >
              {/* 種目名とボタン */}
              <div className="flex items-center gap-2 mb-2">
                <select
                  value={exercise.name}
                  onChange={(e) =>
                    handleChangeExerciseName(exerciseIndex, e.target.value)
                  }
                  className="flex-1 px-2 py-1.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-1 focus:ring-booking-500 focus:border-booking-500 bg-white text-xs font-semibold transition-colors"
                >
                  <option value="">
                    {t("workout.exercises.exercisePlaceholder")}
                  </option>
                  {partExercises.map((partExercise) => (
                    <option key={partExercise.id} value={partExercise.name}>
                      {partExercise.name}
                    </option>
                  ))}
                </select>
                {exercise.id && lastRecords.get(exercise.id) && (
                  <button
                    type="button"
                    onClick={() => handleCopyLastRecord(exerciseIndex)}
                    className="p-1.5 text-booking-600 hover:text-booking-700 bg-booking-50 hover:bg-booking-100 rounded-lg transition-colors border border-booking-200"
                    title={t("workout.exercises.copyLastRecord")}
                  >
                    <svg
                      className="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      strokeWidth={2.5}
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                      />
                    </svg>
                  </button>
                )}
                <button
                  type="button"
                  onClick={() => handleRemoveExercise(exerciseIndex)}
                  className="p-1 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors"
                  disabled={exercises.length === 1}
                >
                  <svg
                    className="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    strokeWidth={2.5}
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                    />
                  </svg>
                </button>
              </div>

              {/* セットリスト（縦並び） */}
              <div className="space-y-1 mb-1.5">
                {exercise.sets.map(
                  (set: ExerciseRow["sets"][number], setIndex: number) => {
                    const previousRecord = exercise.id
                      ? lastRecords.get(exercise.id)
                      : null;
                    const previousSet = previousRecord?.sets?.[setIndex];

                    return (
                      <div
                        key={setIndex}
                        className="flex items-center gap-1 bg-gradient-to-r from-gray-50 to-gray-50/50 rounded-lg p-1.5 border border-gray-100"
                      >
                        <span className="text-xs font-bold text-gray-600 w-4 flex-shrink-0">
                          {setIndex + 1}
                        </span>
                        <input
                          type="number"
                          value={set.weight_kg as any}
                          onChange={(e) =>
                            handleUpdateCell(
                              exerciseIndex,
                              setIndex,
                              "weight_kg",
                              e.target.value,
                            )
                          }
                          className="w-12 px-1 py-1 text-xs font-semibold border border-gray-300 rounded-md text-center focus:outline-none focus:ring-1 focus:ring-booking-500 focus:border-booking-500 bg-white transition-colors [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
                        />
                        <span className="text-[9px] font-medium text-gray-600">
                          {t("workout.exercises.kg")}
                        </span>
                        <span className="text-gray-300 text-xs font-bold">
                          ×
                        </span>
                        <input
                          type="number"
                          value={set.reps as any}
                          onChange={(e) =>
                            handleUpdateCell(
                              exerciseIndex,
                              setIndex,
                              "reps",
                              e.target.value,
                            )
                          }
                          className="w-12 px-1 py-1 text-xs font-semibold border border-gray-300 rounded-md text-center focus:outline-none focus:ring-1 focus:ring-booking-500 focus:border-booking-500 bg-white transition-colors [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
                        />
                        <span className="text-[9px] font-medium text-gray-600">
                          {t("workout.exercises.reps")}
                        </span>
                        {previousSet && (
                          <span className="text-[9px] text-gray-500 ml-1">
                            {t("workout.exercises.previousLabel")} :{" "}
                            {previousSet.weight_kg}
                            {t("workout.exercises.kg")}×{previousSet.reps}
                          </span>
                        )}
                        <div className="ml-auto flex items-center gap-0.5">
                          <button
                            type="button"
                            onClick={() =>
                              handleCopySetBelow(exerciseIndex, setIndex)
                            }
                            className="p-0.5 text-purple-600 hover:text-purple-700 bg-purple-50 hover:bg-purple-100 rounded transition-colors border border-purple-200 disabled:opacity-30 disabled:cursor-not-allowed"
                            disabled={exercise.sets.length >= 5}
                            title={t("workout.exercises.copySetBelow")}
                          >
                            <svg
                              className="w-3.5 h-3.5"
                              fill="none"
                              stroke="currentColor"
                              viewBox="0 0 24 24"
                              strokeWidth={2.5}
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                d="M15 13l-3 3m0 0l-3-3m3 3V8m0 13a9 9 0 110-18 9 9 0 010 18z"
                              />
                            </svg>
                          </button>
                          <button
                            type="button"
                            onClick={() =>
                              handleRemoveSet(exerciseIndex, setIndex)
                            }
                            className="p-0.5 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
                            disabled={exercise.sets.length === 1}
                          >
                            <svg
                              className="w-3.5 h-3.5"
                              fill="none"
                              stroke="currentColor"
                              viewBox="0 0 24 24"
                              strokeWidth={2.5}
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                d="M6 18L18 6M6 6l12 12"
                              />
                            </svg>
                          </button>
                        </div>
                      </div>
                    );
                  },
                )}
              </div>

              {/* セット追加ボタン（5セット未満の場合のみ表示） */}
              {exercise.sets.length < 5 && (
                <div className="flex justify-center mb-1.5">
                  <button
                    type="button"
                    onClick={() => handleAddSet(exerciseIndex)}
                    className="p-1 text-booking-600 hover:text-booking-700 bg-booking-50 hover:bg-booking-100 rounded-lg border border-booking-200 transition-colors"
                  >
                    <svg
                      className="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      strokeWidth={2.5}
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        d="M12 4v16m8-8H4"
                      />
                    </svg>
                  </button>
                </div>
              )}

              {/* メモ */}
              <div>
                <input
                  type="text"
                  value={exercise.sets[0]?.note ?? ""}
                  placeholder={t("workout.exercises.notePlaceholder")}
                  onChange={(e) =>
                    handleUpdateNote(exerciseIndex, e.target.value)
                  }
                  className="w-full px-2 py-1.5 text-xs text-gray-700 border border-gray-300 rounded-lg focus:outline-none focus:ring-1 focus:ring-booking-500 focus:border-booking-500 bg-gray-50 transition-colors placeholder:text-gray-400"
                />
              </div>
            </div>
          ))}

          {/* 種目追加ボタン */}
          <div className="flex justify-center pt-0.5">
            <button
              type="button"
              onClick={handleAddExerciseRow}
              className="flex items-center gap-1.5 px-4 py-2 text-xs text-booking-600 font-bold hover:text-booking-700 bg-booking-50 hover:bg-booking-100 rounded-xl border border-booking-200 hover:border-booking-300 transition-all shadow-sm hover:shadow"
            >
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                strokeWidth={2.5}
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M12 4v16m8-8H4"
                />
              </svg>
              <span>{t("workout.exercises.addExerciseButton")}</span>
            </button>
          </div>
        </div>

        {/* 登録/更新ボタン（モバイル: 下部固定） */}
        <div className="md:hidden fixed bottom-0 left-0 right-0 p-3 bg-white border-t border-gray-200 shadow-lg z-10">
          <button
            onClick={onSubmit}
            className="w-full py-2.5 text-sm rounded-xl bg-booking-600 text-white hover:bg-booking-700 active:bg-booking-800 transition-colors font-bold shadow-sm"
          >
            {isUpdate
              ? t("workout.exercises.updateButtonMobile")
              : t("workout.exercises.registerButtonMobile")}
          </button>
        </div>
      </>
    );
  },
);

MobileView.displayName = "MobileView";

export default MobileView;
