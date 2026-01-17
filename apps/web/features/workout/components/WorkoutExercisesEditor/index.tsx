"use client";
import { useState, useEffect, useCallback } from "react";
import { useTranslation } from "react-i18next";
import type { ExerciseRow } from "@/features/workout/lib/utils";
import type { WorkoutPartDTO, ExerciseDTO } from "@/types/workout";
import {
  updateExerciseCell,
  updateExerciseNote,
  createEmptyExerciseRow,
} from "@/features/workout/lib/utils";
import ExerciseManageModal from "../ExerciseManageModal";
import { useIsMobile } from "@/features/workout/hooks/useIsMobile";
import { useMobileExerciseAdjustment } from "@/features/workout/hooks/useMobileExerciseAdjustment";
import { useWorkoutContext } from "@/app/workout/_providers/WorkoutProvider";
import { useFormContext } from "react-hook-form";
import type { WorkoutFormDTO } from "@/types/workout";
import MobileView from "./MobileView";
import DesktopView from "./DesktopView";

const WorkoutExercisesEditor = () => {
  const { t, i18n } = useTranslation("common");
  const {
    availableParts: workoutParts,
    selectedPartId,
    setSelectedPartId: onPartChange,
    isUpdate,
    handleSubmit: onSubmit,
    refetchWorkoutParts: onRefetchParts,
    fetchLastRecord: onFetchLastRecord,
    dateStr: dataKey,
  } = useWorkoutContext();
  const form = useFormContext<WorkoutFormDTO>();

  const allExercises = form.watch("exercises") || [];
  const onChangeExercises = useCallback(
    (exercises: ExerciseRow[]) => {
      form.setValue("exercises", exercises, { shouldDirty: true });
    },
    [form],
  );

  const displayedExercises = selectedPartId
    ? allExercises.filter((ex) => ex.workout_part_id === selectedPartId)
    : [];

  const exercises =
    displayedExercises.length > 0
      ? displayedExercises
      : [createEmptyExerciseRow(1)];

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [partExercises, setPartExercises] = useState<
    Array<{ id: number; name: string; workout_part_id: number | null }>
  >([]);
  const [lastRecords, setLastRecords] = useState<Map<number, ExerciseDTO>>(
    new Map(),
  );

  const isMobile = useIsMobile();

  useMobileExerciseAdjustment(
    isMobile,
    allExercises,
    onChangeExercises,
    dataKey,
  );

  useEffect(() => {
    if (!selectedPartId) {
      setPartExercises([]);
      return;
    }
    const part = workoutParts.find(
      (workoutPart) => workoutPart.id === selectedPartId,
    );
    setPartExercises(part?.exercises || []);
  }, [selectedPartId, workoutParts]);

  useEffect(() => {
    const fetchPreviousRecords = async () => {
      const newLastRecords = new Map<number, ExerciseDTO>();

      for (const exercise of exercises) {
        if (exercise.id) {
          try {
            const lastRecord = await onFetchLastRecord(exercise.id);
            if (lastRecord) {
              newLastRecords.set(exercise.id, lastRecord);
            }
          } catch (err) {
            console.error(
              `Failed to fetch last record for exercise ${exercise.id}:`,
              err,
            );
          }
        }
      }

      setLastRecords(newLastRecords);
    };

    fetchPreviousRecords();
  }, [exercises.map((e) => e.id).join(","), onFetchLastRecord]);

  const updateDisplayedExercise = useCallback(
    (updatedDisplayed: ExerciseRow[]) => {
      const otherExercises = allExercises.filter(
        (ex) => ex.workout_part_id !== selectedPartId,
      );
      onChangeExercises([...otherExercises, ...updatedDisplayed]);
    },
    [allExercises, selectedPartId, onChangeExercises],
  );

  const handleUpdateCell = useCallback(
    (
      exerciseIndex: number,
      setIndex: number,
      key: "weight_kg" | "reps",
      value: string,
    ) => {
      const updated = updateExerciseCell(
        exercises,
        exerciseIndex,
        setIndex,
        key,
        value,
      );
      updateDisplayedExercise(updated);
    },
    [exercises, updateDisplayedExercise],
  );

  const handleUpdateNote = useCallback(
    (exerciseIndex: number, note: string) => {
      const updated = updateExerciseNote(exercises, exerciseIndex, note);
      updateDisplayedExercise(updated);
    },
    [exercises, updateDisplayedExercise],
  );

  const handleRemoveSet = useCallback(
    (exerciseIndex: number, setIndex: number) => {
      if (exercises[exerciseIndex].sets.length === 1) return;
      const updatedExercises = structuredClone(exercises);
      updatedExercises[exerciseIndex].sets.splice(setIndex, 1);
      updateDisplayedExercise(updatedExercises);
    },
    [exercises, updateDisplayedExercise],
  );

  const handleChangeExerciseName = useCallback(
    (exerciseIndex: number, exerciseName: string) => {
      const selectedExercise = partExercises.find(
        (exercise) => exercise.name === exerciseName,
      );
      const updatedExercises = structuredClone(exercises);
      updatedExercises[exerciseIndex].name = exerciseName;
      updatedExercises[exerciseIndex].id = selectedExercise?.id || null;
      updatedExercises[exerciseIndex].workout_part_id =
        selectedExercise?.workout_part_id || selectedPartId || null;
      updateDisplayedExercise(updatedExercises);
    },
    [exercises, partExercises, selectedPartId, updateDisplayedExercise],
  );

  const handleAddExerciseRow = useCallback(() => {
    const newExercise = createEmptyExerciseRow(1);
    newExercise.workout_part_id = selectedPartId;
    updateDisplayedExercise([...exercises, newExercise]);
  }, [exercises, selectedPartId, updateDisplayedExercise]);

  const handleAddSet = useCallback(
    (exerciseIndex: number) => {
      const updatedExercises = structuredClone(exercises);
      const newSetNumber = updatedExercises[exerciseIndex].sets.length + 1;
      updatedExercises[exerciseIndex].sets.push({
        set_number: newSetNumber,
        weight_kg: 0,
        reps: 0,
        note: null,
      });
      updateDisplayedExercise(updatedExercises);
    },
    [exercises, updateDisplayedExercise],
  );

  const handleRemoveExercise = useCallback(
    (exerciseIndex: number) => {
      if (exercises.length === 1) return;
      updateDisplayedExercise(
        exercises.filter((_, index) => index !== exerciseIndex),
      );
    },
    [exercises, updateDisplayedExercise],
  );

  const handleCopyLastRecord = useCallback(
    (exerciseIndex: number) => {
      const currentExercise = exercises[exerciseIndex];
      if (!currentExercise.id) return;

      const previousRecord = lastRecords.get(currentExercise.id);
      if (
        !previousRecord ||
        !previousRecord.sets ||
        previousRecord.sets.length === 0
      )
        return;

      const updatedExercises = structuredClone(exercises);
      const previousSets = previousRecord.sets;
      updatedExercises[exerciseIndex].sets = previousSets.map(
        (previousSet: ExerciseDTO["sets"][number], index: number) => ({
          id: null,
          set_number: index + 1,
          weight_kg: String(previousSet.weight_kg || ""),
          reps: String(previousSet.reps || ""),
          note: previousSet.note || null,
        }),
      );

      updateDisplayedExercise(updatedExercises);
    },
    [exercises, lastRecords, updateDisplayedExercise],
  );

  const handleCopySetBelow = useCallback(
    (exerciseIndex: number, setIndex: number) => {
      const currentExercise = exercises[exerciseIndex];
      if (!currentExercise.sets || currentExercise.sets.length >= 5) return;

      const currentSet = currentExercise.sets[setIndex];
      const updatedExercises = structuredClone(exercises);

      const copiedSet = {
        id: null,
        set_number: setIndex + 2,
        weight_kg: currentSet.weight_kg || 0,
        reps: currentSet.reps || 0,
        note: null,
      };

      updatedExercises[exerciseIndex].sets.splice(setIndex + 1, 0, copiedSet);
      updatedExercises[exerciseIndex].sets = updatedExercises[
        exerciseIndex
      ].sets.map((set: ExerciseRow["sets"][number], index: number) => ({
        ...set,
        set_number: index + 1,
      }));

      updateDisplayedExercise(updatedExercises);
    },
    [exercises, updateDisplayedExercise],
  );

  const handlePartChange = useCallback(
    (idStr: string) => {
      if (!idStr) {
        onPartChange(null);
        return;
      }
      const numId = Number(idStr);
      onPartChange(numId);
    },
    [onPartChange],
  );

  const handleSuccess = useCallback(() => {
    onRefetchParts();
  }, [onRefetchParts]);

  const commonProps = {
    t,
    i18n,
    exercises,
    workoutParts,
    selectedPartId,
    partExercises,
    lastRecords,
    isUpdate,
    isModalOpen,
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
    handleSuccess,
  };

  return (
    <>
      <ExerciseManageModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        workoutParts={workoutParts}
        onSuccess={handleSuccess}
      />
      {isMobile ? (
        <MobileView {...commonProps} />
      ) : (
        <DesktopView {...commonProps} />
      )}
    </>
  );
};

export default WorkoutExercisesEditor;
