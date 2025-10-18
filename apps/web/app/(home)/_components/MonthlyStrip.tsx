"use client";
import React from "react";

export default function MonthlyStrip({ year, month, selectedDay, onSelectDay }: { year: number; month: number; selectedDay: number; onSelectDay: (d: number) => void }) {
  const daysOfWeek = ["日", "月", "火", "水", "木", "金", "土"];
  const now = new Date();
  const daysInMonth = new Date(year, month, 0).getDate();
  const days = Array.from({ length: daysInMonth }, (_, i) => i + 1);

  return (
    <div className="bg-white rounded-lg shadow mb-6 p-6 overflow-x-auto">
      <div className="flex gap-2 min-w-max">
        {days.map((day) => {
          const date = new Date(year, month - 1, day);
          const label = daysOfWeek[date.getDay()];
          const active = selectedDay === day;
          return (
            <button
              key={day}
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
