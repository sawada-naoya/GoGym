"use client";
import { memo, Fragment } from "react";
import type { ExerciseRow } from "@/features/workout/lib/utils";
import type { WorkoutPartDTO, ExerciseDTO } from "@/types/workout";
import { getLocalizedPartName } from "@/features/workout/lib/utils";

type DesktopViewProps = {
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

const DesktopView = memo<DesktopViewProps>(
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
    handleUpdateCell,
    handleAddSet,
    handleUpdateNote,
    handleAddExerciseRow,
    onSubmit,
  }) => {
    return (
      <div className="hidden md:block bg-white rounded-lg shadow mb-4 md:mb-6 p-3 md:p-6">
        {/* ヘッダー: 部位選択と種目作成ボタン、登録/更新ボタン */}
        <div className="flex items-center justify-between gap-2 mb-6">
          <div className="flex items-center gap-2">
            <label className="text-base font-medium text-gray-700">
              {t("workout.exercises.partLabel")}
            </label>
            <select
              value={selectedPartId?.toString() ?? ""}
              onChange={(e) => handlePartChange(e.target.value)}
              className="w-40 px-3 py-2.5 text-base border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-booking-500 bg-white"
            >
              <option value="">{t("workout.exercises.partPlaceholder")}</option>
              {workoutParts.map((part) => (
                <option key={part.id} value={part.id.toString()}>
                  {getLocalizedPartName(part, i18n.language)}
                </option>
              ))}
            </select>
            <button
              type="button"
              onClick={() => setIsModalOpen(true)}
              className="px-4 py-2.5 rounded-md text-booking-600 hover:bg-booking-50 transition-colors whitespace-nowrap text-base border border-gray-200 hover:border-booking-300 font-medium"
            >
              {t("workout.exercises.createExerciseButton")}
            </button>
          </div>

          <button
            onClick={onSubmit}
            className="px-8 py-2.5 text-base rounded-md bg-booking-600 text-white hover:bg-booking-700 transition-colors font-medium"
          >
            {isUpdate
              ? t("workout.exercises.updateButton")
              : t("workout.exercises.registerButton")}
          </button>
        </div>

        {/* テーブル表示 */}
        <div className="overflow-x-auto">
          <table className="w-full border-collapse">
            <thead>
              <tr className="border-b-2 border-gray-300">
                <th className="px-4 py-3 text-left font-medium text-gray-700 border-r border-gray-300 min-w-[240px]">
                  {t("workout.exercises.exercisePlaceholder")}
                </th>
                {[1, 2, 3, 4, 5].map((setNumber) => (
                  <th
                    key={setNumber}
                    colSpan={2}
                    className="px-4 py-3 text-center font-medium text-gray-700 border-r border-gray-300"
                  >
                    {setNumber}
                    {t("workout.exercises.setLabel")}
                  </th>
                ))}
                <th className="px-4 py-3 text-center font-medium text-gray-700 w-12"></th>
              </tr>
              <tr className="border-b border-gray-300 bg-gray-50">
                <th className="px-4 py-2 border-r border-gray-300" />
                {[1, 2, 3, 4, 5].map((setNumber) => (
                  <Fragment key={setNumber}>
                    <th className="px-2 py-2 text-center text-sm font-medium text-gray-600 border-r border-gray-200">
                      {t("workout.exercises.weightLabel")}
                    </th>
                    <th className="px-2 py-2 text-center text-sm font-medium text-gray-600 border-r border-gray-300">
                      {t("workout.exercises.repsLabel")}
                    </th>
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
                        <select
                          value={exercise.name}
                          onChange={(e) =>
                            handleChangeExerciseName(
                              exerciseIndex,
                              e.target.value,
                            )
                          }
                          className="flex-1 px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-booking-500 bg-white text-sm truncate"
                        >
                          <option value="">
                            {t("workout.exercises.exercisePlaceholder")}
                          </option>
                          {partExercises.map((partExercise) => (
                            <option
                              key={partExercise.id}
                              value={partExercise.name}
                            >
                              {partExercise.name}
                            </option>
                          ))}
                        </select>
                        {exercise.id && lastRecords.get(exercise.id) && (
                          <button
                            type="button"
                            onClick={() => handleCopyLastRecord(exerciseIndex)}
                            className="p-1 text-booking-600 hover:text-booking-700 hover:bg-booking-50 rounded transition-colors flex-shrink-0"
                            title={t("workout.exercises.copyLastRecord")}
                          >
                            <svg
                              className="w-4 h-4"
                              fill="none"
                              stroke="currentColor"
                              viewBox="0 0 24 24"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                              />
                            </svg>
                          </button>
                        )}
                      </div>
                    </td>

                    {[0, 1, 2, 3, 4].map((setIndex) => {
                      const set = exercise.sets[setIndex];
                      const previousRecord = exercise.id
                        ? lastRecords.get(exercise.id)
                        : null;
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
                                    onChange={(e) =>
                                      handleUpdateCell(
                                        exerciseIndex,
                                        setIndex,
                                        "weight_kg",
                                        e.target.value,
                                      )
                                    }
                                    className="w-14 px-2 py-1 border border-gray-300 rounded text-center focus:outline-none focus:ring-1 focus:ring-booking-500 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
                                  />
                                  <span className="text-xs text-gray-600">
                                    {t("workout.exercises.kg")}
                                  </span>
                                </div>
                                {previousSet && (
                                  <span className="text-[10px] text-gray-500 text-center">
                                    {t("workout.exercises.previousLabel")}:{" "}
                                    {previousSet.weight_kg}
                                    {t("workout.exercises.kg")}
                                  </span>
                                )}
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
                                    onChange={(e) =>
                                      handleUpdateCell(
                                        exerciseIndex,
                                        setIndex,
                                        "reps",
                                        e.target.value,
                                      )
                                    }
                                    className="w-14 px-2 py-1 border border-gray-300 rounded text-center focus:outline-none focus:ring-1 focus:ring-booking-500 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
                                  />
                                  <span className="text-xs text-gray-600">
                                    {t("workout.exercises.rep")}
                                  </span>
                                </div>
                                {previousSet && (
                                  <span className="text-[10px] text-gray-500 text-center">
                                    {t("workout.exercises.previousLabel")}:{" "}
                                    {previousSet.reps}
                                    {t("workout.exercises.reps")}
                                  </span>
                                )}
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
                        <button
                          type="button"
                          onClick={() => handleAddSet(exerciseIndex)}
                          className="p-1 text-booking-600 hover:text-booking-700 hover:bg-booking-50 rounded transition-colors"
                        >
                          <svg
                            className="w-5 h-5"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M12 4v16m8-8H4"
                            />
                          </svg>
                        </button>
                      )}
                    </td>
                  </tr>

                  <tr className="border-b border-gray-300 hover:bg-gray-50">
                    <td colSpan={12} className="px-4 py-2">
                      <div className="flex items-center gap-2">
                        <input
                          type="text"
                          value={exercise.sets[0]?.note ?? ""}
                          placeholder={t("workout.exercises.notePlaceholder")}
                          onChange={(e) =>
                            handleUpdateNote(exerciseIndex, e.target.value)
                          }
                          className="flex-1 px-3 py-2 text-sm text-gray-700 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-booking-500 truncate"
                        />
                      </div>
                    </td>
                  </tr>
                </Fragment>
              ))}

              <tr className="h-24 border-b border-gray-200">
                <td
                  colSpan={12}
                  className="px-4 py-3 text-center text-gray-400"
                >
                  <button
                    className="inline-flex items-center justify-center w-10 h-10 text-booking-600 hover:text-booking-700 hover:bg-booking-50 rounded-full transition-colors"
                    onClick={handleAddExerciseRow}
                  >
                    <svg
                      className="w-6 h-6"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M12 4v16m8-8H4"
                      />
                    </svg>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    );
  },
);

DesktopView.displayName = "DesktopView";

export default DesktopView;
