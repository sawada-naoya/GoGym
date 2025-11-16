"use client";
import React, { useEffect, useRef } from "react";

export default function MonthlyStrip({ year, month, selectedDay, onSelectDay }: { year: number; month: number; selectedDay: number; onSelectDay: (d: number) => void }) {
  const daysOfWeek = ["日", "月", "火", "水", "木", "金", "土"];
  const daysInMonth = new Date(year, month, 0).getDate();
  const days = Array.from({ length: daysInMonth }, (_, i) => i + 1);

  const containerRef = useRef<HTMLDivElement>(null);
  const selectedButtonRef = useRef<HTMLButtonElement>(null);

  // 初期表示時に選択された日付が中央に来るようにスクロール
  useEffect(() => {
    if (containerRef.current && selectedButtonRef.current) {
      const container = containerRef.current;
      const button = selectedButtonRef.current;

      const containerWidth = container.clientWidth;
      const buttonLeft = button.offsetLeft;
      const buttonWidth = button.clientWidth;

      // ボタンの中心位置 - コンテナの中心位置 = スクロール位置
      const scrollPosition = buttonLeft - containerWidth / 2 + buttonWidth / 2;

      container.scrollLeft = scrollPosition;
    }
  }, []); // 初回のみ実行

  return (
    <div ref={containerRef} className="overflow-x-auto">
      <div className="flex gap-2 min-w-max">
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
              className={`flex flex-col items-center justify-center w-16 h-16 rounded-lg border-2 transition-colors
                ${active ? "border-booking-600 bg-booking-50 text-booking-700" : "border-gray-200 hover:border-booking-400 hover:bg-gray-50"}`}
            >
              <span className="text-lg font-bold">{day}</span>
              <span className="text-xs text-gray-600">({label})</span>
            </button>
          );
        })}
      </div>
    </div>
  );
}
