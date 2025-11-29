"use client";
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
    <div className="bg-white rounded-lg shadow mb-6 p-6">
      {/* 年月選択 */}
      <div className="flex items-center justify-center gap-4 mb-6">
        <select
          value={selectedYear}
          onChange={(e) => {
            const year = Number(e.target.value);
            onYearChange(year);
            handleYearMonthChange(year, selectedMonth);
          }}
          className="w-28 text-center border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-booking-500 bg-white"
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
        <span className="text-lg font-medium">年</span>

        <select
          value={selectedMonth}
          onChange={(e) => {
            const month = Number(e.target.value);
            onMonthChange(month);
            handleYearMonthChange(selectedYear, month);
          }}
          className="w-20 text-center border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-booking-500 bg-white"
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
        <span className="text-lg font-medium">月</span>
      </div>

      {/* 日付選択 */}
      <MonthlyStrip year={selectedYear} month={selectedMonth} selectedDay={selectedDay} onSelectDay={handleDayChange} />

      {/* 区切り線 */}
      <div className="my-6 border-t border-gray-200"></div>

      {/* トレーニング詳細（時刻・場所・コンディション） */}
      <div className="flex flex-wrap items-start justify-start gap-3 text-left">
        {/* 時間 */}
        <div className="flex items-center gap-1">
          <label className="mr-2 text-sm font-medium text-gray-700">時間</label>
          <input type="time" className="w-24 px-2 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" value={form.watch("started_at") ?? ""} onChange={(e) => form.setValue("started_at", e.target.value || null, { shouldDirty: true })} />
          <span className="text-gray-500">〜</span>
          <input type="time" className="w-24 px-2 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" value={form.watch("ended_at") ?? ""} onChange={(e) => form.setValue("ended_at", e.target.value || null, { shouldDirty: true })} />
        </div>

        {/* 場所 */}
        <div className="flex items-center gap-2">
          <label className="text-sm font-medium text-gray-700">場所</label>
          <input type="text" value={form.watch("place") ?? ""} onChange={(e) => form.setValue("place", e.target.value, { shouldDirty: true })} placeholder="ジム名など" className="w-80 px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" />
        </div>

        {/* コンディション（1〜5） */}
        <div className="flex items-center gap-2">
          <label className="text-sm font-medium text-gray-700">コンディション</label>
          <div className="flex gap-1">
            {[1, 2, 3, 4, 5].map((n) => (
              <button key={n} type="button" onClick={() => form.setValue("condition_level", n as 1 | 2 | 3 | 4 | 5, { shouldDirty: true })} className={`w-10 h-10 rounded border text-sm font-medium transition-colors ${form.watch("condition_level") === n ? "bg-booking-600 text-white border-booking-600" : "border-gray-300 text-gray-700 hover:bg-gray-100"}`}>
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
