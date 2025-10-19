"use client";
import React from "react";
import type { WorkoutFormDTO } from "./types";

type SessionMeta = Pick<WorkoutFormDTO, "startedAt" | "endedAt" | "place" | "conditionLevel" | "workoutPart">;

type Props = {
  value: SessionMeta;
  onChange: (next: SessionMeta) => void;
};

const WorkoutSessionMetaEditor: React.FC<Props> = ({ value, onChange }) => {
  const setTime = (key: "startedAt" | "endedAt", hhmm: string) => {
    onChange({ ...value, [key]: hhmm || null });
  };

  const setPlace = (place: string) => {
    onChange({ ...value, place });
  };

  const setCondition = (n: 1 | 2 | 3 | 4 | 5) => {
    onChange({ ...value, conditionLevel: n });
  };

  // 部位は自由入力。空なら null、文字があれば custom 扱いで保持。
  const setPartName = (name: string) => {
    const trimmed = name.trim();
    if (!trimmed) {
      onChange({
        ...value,
        workoutPart: { id: null, name: null, source: null },
      });
      return;
    }
    onChange({
      ...value,
      workoutPart: { id: null, name: trimmed, source: "custom" },
    });
  };

  return (
    <div className="flex flex-wrap items-start justify-start gap-3 text-left">
      {/* 時間 */}
        <div className="flex items-center gap-1">
          <label className="mr-2 text-sm font-medium text-gray-700">時間</label>
          <input type="time" className="w-24 px-2 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" value={value.startedAt ?? ""} onChange={(e) => setTime("startedAt", e.target.value)} />
          <span className="text-gray-500">〜</span>
          <input type="time" className="w-24 px-2 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" value={value.endedAt ?? ""} onChange={(e) => setTime("endedAt", e.target.value)} />
        </div>

        {/* 場所 */}
        <div className="flex items-center gap-2">
          <label className="text-sm font-medium text-gray-700">場所</label>
          <input type="text" value={value.place ?? ""} onChange={(e) => setPlace(e.target.value)} placeholder="自宅 / ジム名など" className="w-48 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" />
        </div>

        {/* 部位（自由入力） */}
        <div className="flex items-center gap-2">
          <label className="text-sm font-medium text-gray-700">部位</label>
          <input type="text" value={value.workoutPart?.name ?? ""} onChange={(e) => setPartName(e.target.value)} placeholder="例: 胸 / 脚 / 背中" className="w-40 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" />
        </div>

        {/* コンディション（1〜5） */}
        <div className="flex items-center gap-2">
          <label className="text-sm font-medium text-gray-700">コンディション</label>
          <div className="flex gap-1">
            {[1, 2, 3, 4, 5].map((n) => (
              <button key={n} type="button" onClick={() => setCondition(n as 1 | 2 | 3 | 4 | 5)} className={`w-10 h-10 rounded border text-sm font-medium transition-colors ${value.conditionLevel === n ? "bg-booking-600 text-white border-booking-600" : "border-gray-300 text-gray-700 hover:bg-gray-100"}`}>
                {n}
              </button>
            ))}
          </div>
        </div>
    </div>
  );
};

export default WorkoutSessionMetaEditor;
