"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { useCallback, useMemo } from "react";
import { formatDate, parseDateString } from "@/utils/time";

/**
 * ワークアウト日付管理hook
 *
 * URLをSingle Source of Truthとして日付を管理
 * - searchParams.date から日付を取得
 * - 日付変更はURL更新のみ（state不要）
 * - ブラウザバック/フォワード対応
 */
export const useWorkoutDate = () => {
  const router = useRouter();
  const searchParams = useSearchParams();

  // URLから日付を取得
  const dateStr = searchParams.get("date");

  const { year, month, day } = useMemo(() => {
    if (dateStr) {
      return parseDateString(dateStr);
    }

    const now = new Date();
    return {
      year: now.getFullYear(),
      month: now.getMonth() + 1,
      day: now.getDate(),
    };
  }, [dateStr]);

  // 日付変更（URL更新のみ）
  const setDate = useCallback(
    (newYear: number, newMonth: number, newDay: number) => {
      const newDateStr = formatDate(newYear, newMonth, newDay);
      // 現在のURLと同じなら何もしない
      if (newDateStr === dateStr) return;

      router.push(`/workout?date=${newDateStr}`, { scroll: false });
    },
    [router, dateStr]
  );

  const setYear = useCallback(
    (newYear: number) => {
      setDate(newYear, month, day);
    },
    [setDate, month, day]
  );

  const setMonth = useCallback(
    (newMonth: number) => {
      setDate(year, newMonth, day);
    },
    [setDate, year, day]
  );

  const setDay = useCallback(
    (newDay: number) => {
      setDate(year, month, newDay);
    },
    [setDate, year, month]
  );

  return {
    year,
    month,
    day,
    dateStr: dateStr || formatDate(year, month, day),
    setDate,
    setYear,
    setMonth,
    setDay,
  };
};
