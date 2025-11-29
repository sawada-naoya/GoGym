"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { UseFormReturn } from "react-hook-form";
import MonthlyStrip from "./MonthlyStrip";
import { WorkoutFormDTO } from "../_lib/types";
import { formatDate, validateDay } from "@/lib/time";

type Props = {
  form: UseFormReturn<WorkoutFormDTO>;
  selectedYear: number;
  selectedMonth: number;
  selectedDay: number;
  onYearChange: (year: number) => void;
  onMonthChange: (month: number) => void;
  onDayChange: (day: number) => void;
};

const WorkoutMetadataEditor = ({ form, selectedYear, selectedMonth, selectedDay, onYearChange, onMonthChange, onDayChange }: Props) => {
  const router = useRouter();
  const [isMetadataOpen, setIsMetadataOpen] = useState(false);

  const handleYearMonthChange = (year: number, month: number) => {
    const validDay = validateDay(year, month, selectedDay);
    const formattedDate = formatDate(year, month, validDay);
    router.push(`/workout?date=${formattedDate}`);
  };

  const handleDayChange = (day: number) => {
    onDayChange(day);
    const formattedDate = formatDate(selectedYear, selectedMonth, day);
    router.push(`/workout?date=${formattedDate}`);
  };

  return (
    <div className="bg-white rounded-lg shadow mb-4 md:mb-6 p-3 md:p-6">
      {/* 年月選択 */}
      <div className="flex items-center justify-center gap-2 md:gap-4 mb-3 md:mb-6">
        <select
          value={selectedYear}
          onChange={(e) => {
            const year = Number(e.target.value);
            onYearChange(year);
            handleYearMonthChange(year, selectedMonth);
          }}
          className="w-20 md:w-28 text-center border border-gray-300 rounded-md px-2 md:px-3 py-1.5 md:py-2 text-sm md:text-base focus:outline-none focus:ring-2 focus:ring-booking-500 bg-white"
        >
          {Array.from({ length: 10 }, (_, i) => {
            const year = new Date().getFullYear() - 5 + i;
            return (
              <option key={year} value={year}>
                {year}
              </option>
            );
          })}
        </select>
        <span className="text-sm md:text-lg font-medium">年</span>

        <select
          value={selectedMonth}
          onChange={(e) => {
            const month = Number(e.target.value);
            onMonthChange(month);
            handleYearMonthChange(selectedYear, month);
          }}
          className="w-16 md:w-20 text-center border border-gray-300 rounded-md px-2 md:px-3 py-1.5 md:py-2 text-sm md:text-base focus:outline-none focus:ring-2 focus:ring-booking-500 bg-white"
        >
          {Array.from({ length: 12 }, (_, i) => {
            const month = i + 1;
            return (
              <option key={month} value={month}>
                {month}
              </option>
            );
          })}
        </select>
        <span className="text-sm md:text-lg font-medium">月</span>
      </div>

      {/* 日付選択 */}
      <MonthlyStrip year={selectedYear} month={selectedMonth} selectedDay={selectedDay} onSelectDay={handleDayChange} />

      {/* モバイル: 折りたたみトグル */}
      <div className="mt-2 md:hidden">
        <button type="button" onClick={() => setIsMetadataOpen(!isMetadataOpen)} className="w-full flex items-center justify-between px-3 py-1.5 text-xs font-medium text-gray-700 bg-gray-50 rounded-md hover:bg-gray-100 transition-colors">
          <span>詳細情報</span>
          <svg className={`w-4 h-4 transition-transform ${isMetadataOpen ? "rotate-180" : ""}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
          </svg>
        </button>
      </div>

      {/* 区切り線（デスクトップのみ） */}
      <div className="hidden md:block my-6 border-t border-gray-200"></div>

      {/* トレーニング詳細（時刻・場所・コンディション） */}
      <div className={`flex flex-col md:flex-row md:flex-wrap items-start justify-start gap-2 md:gap-3 text-left ${isMetadataOpen ? "mt-2" : "hidden"} md:flex`}>
        {/* 時間 */}
        <div className="flex items-center gap-1 w-full md:w-auto">
          <label className="text-xs md:text-sm font-medium text-gray-700 min-w-[44px] md:min-w-[60px]">時間</label>
          <input type="time" className="w-24 md:w-28 px-1.5 md:px-2 py-1.5 md:py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" value={form.watch("started_at") ?? ""} onChange={(e) => form.setValue("started_at", e.target.value || null, { shouldDirty: true })} />
          <span className="text-xs text-gray-500">〜</span>
          <input type="time" className="w-24 md:w-28 px-1.5 md:px-2 py-1.5 md:py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" value={form.watch("ended_at") ?? ""} onChange={(e) => form.setValue("ended_at", e.target.value || null, { shouldDirty: true })} />
        </div>

        {/* 場所 */}
        <div className="flex items-center gap-2 w-full md:w-auto">
          <label className="text-xs md:text-sm font-medium text-gray-700 min-w-[44px] md:min-w-[60px]">場所</label>
          <input type="text" value={form.watch("place") ?? ""} onChange={(e) => form.setValue("place", e.target.value, { shouldDirty: true })} placeholder="ジム名など" className="flex-1 md:w-80 px-2 md:px-3 py-1.5 md:py-2 text-xs md:text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" />
        </div>

        {/* コンディション（1〜5） */}
        <div className="flex items-center gap-2 w-full md:w-auto">
          <label className="text-xs md:text-sm font-medium text-gray-700 min-w-[44px] md:min-w-[60px]">体調</label>
          <div className="flex gap-1">
            {[1, 2, 3, 4, 5].map((n) => (
              <button
                key={n}
                type="button"
                onClick={() => form.setValue("condition_level", n as 1 | 2 | 3 | 4 | 5, { shouldDirty: true })}
                className={`w-8 h-8 md:w-10 md:h-10 rounded border text-xs md:text-sm font-medium transition-colors ${form.watch("condition_level") === n ? "bg-booking-600 text-white border-booking-600" : "border-gray-300 text-gray-700 hover:bg-gray-100"}`}
              >
                {n}
              </button>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default WorkoutMetadataEditor;
