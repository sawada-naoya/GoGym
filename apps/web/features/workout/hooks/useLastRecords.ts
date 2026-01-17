"use client";

import { useState, useCallback } from "react";
import { getLastExerciseRecord } from "@/features/workout/actions";
import type { ExerciseDTO } from "@/types/workout";

/**
 * 種目の前回記録を管理するhook
 *
 * 責務:
 * - 種目IDから前回記録を取得
 * - 取得済みデータをキャッシュ（Map管理）
 */
export const useLastRecords = () => {
  const [lastRecords, setLastRecords] = useState<Map<number, ExerciseDTO>>(
    new Map(),
  );

  /**
   * 特定の種目IDの前回記録を取得
   */
  const fetchLastRecord = useCallback(
    async (exerciseId: number): Promise<ExerciseDTO | null> => {
      if (!exerciseId) return null;

      // 既にキャッシュにある場合はそれを返す
      const cached = lastRecords.get(exerciseId);
      if (cached) return cached;

      try {
        const result = await getLastExerciseRecord(exerciseId);
        if (result.success && result.data) {
          const data = result.data; // 変数に代入で型を確定
          // キャッシュに保存
          setLastRecords((prev) => {
            const next = new Map(prev);
            next.set(exerciseId, data);
            return next;
          });
          return data;
        }
        return null;
      } catch {
        return null;
      }
    },
    [lastRecords],
  );
  /**
   * キャッシュをクリア
   */
  const clearCache = useCallback(() => {
    setLastRecords(new Map());
  }, []);

  return {
    fetchLastRecord,
    clearCache,
    lastRecords,
  };
};
