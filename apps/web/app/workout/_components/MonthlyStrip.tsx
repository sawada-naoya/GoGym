"use client";
import { useRef } from "react";
import { useScrollToCenter } from "../_hooks";

export default function MonthlyStrip({ year, month, selectedDay, onSelectDay }: { year: number; month: number; selectedDay: number; onSelectDay: (d: number) => void }) {
  const daysOfWeek = ["日", "月", "火", "水", "木", "金", "土"];
  const daysInMonth = new Date(year, month, 0).getDate();
  const days = Array.from({ length: daysInMonth }, (_, i) => i + 1);

  const containerRef = useRef<HTMLDivElement>(null);
  const selectedButtonRef = useRef<HTMLButtonElement>(null);

  // カスタムフック: 選択された日付を中央にスクロール
  useScrollToCenter(containerRef, selectedButtonRef);

  return (
    <div ref={containerRef} className="overflow-x-auto">
      <div className="flex gap-1 md:gap-2 min-w-max">
        {days.map((day) => {
          const date = new Date(year, month - 1, day);
          const label = daysOfWeek[date.getDay()];
          const active = selectedDay === day;
          return (
            <button
              key={day}
              ref={active ? selectedButtonRef : null}
              onClick={() => onSelectDay(day)}
              aria-pressed={active}
              className={`flex flex-col items-center justify-center w-11 h-11 md:w-14 md:h-14 rounded-md md:rounded-lg border transition-colors
                ${active ? "border-booking-600 bg-booking-50 text-booking-700 font-semibold" : "border-gray-200 hover:border-booking-400 hover:bg-gray-50"}`}
            >
              <span className="text-sm md:text-base font-bold">{day}</span>
              <span className="text-[10px] md:text-xs text-gray-500">({label})</span>
            </button>
          );
        })}
      </div>
    </div>
  );
}
