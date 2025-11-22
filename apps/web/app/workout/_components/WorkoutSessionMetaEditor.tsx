"use client";
import React from "react";
import type { WorkoutFormDTO, WorkoutPartDTO } from "../_lib/types";

type SessionMeta = Pick<WorkoutFormDTO, "started_at" | "ended_at" | "place" | "condition_level" | "workout_part">;

type Props = {
  value: SessionMeta;
  availableParts: WorkoutPartDTO[];
  onChange: (next: SessionMeta) => void;
};

const WorkoutSessionMetaEditor: React.FC<Props> = ({ value, availableParts, onChange }) => {
  const setTime = (key: "started_at" | "ended_at", hhmm: string) => {
    onChange({ ...value, [key]: hhmm || null });
  };

  const setPlace = (place: string) => {
    onChange({ ...value, place });
  };

  const setCondition = (n: 1 | 2 | 3 | 4 | 5) => {
    onChange({ ...value, condition_level: n });
  };

  // プルダウンで選択された部位をセット
  const setPartById = (idStr: string) => {
    if (!idStr) {
      // 未選択
      onChange({
        ...value,
        workout_part: { id: null, name: null, source: null },
      });
      return;
    }
    const numId = Number(idStr);
    const selectedPart = availableParts.find((p) => {
      return p.id === numId;
    });

    if (selectedPart) {
      onChange({
        ...value,
        workout_part: {
          id: selectedPart.id,
          name: selectedPart.name,
          source: selectedPart.is_default ? "preset" : "custom",
        },
      });
    }
  };

  return (
    <div className="flex flex-wrap items-start justify-start gap-3 text-left">
      {/* 時間 */}
      <div className="flex items-center gap-1">
        <label className="mr-2 text-sm font-medium text-gray-700">時間</label>
        <input type="time" className="w-24 px-2 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" value={value.started_at ?? ""} onChange={(e) => setTime("started_at", e.target.value)} />
        <span className="text-gray-500">〜</span>
        <input type="time" className="w-24 px-2 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" value={value.ended_at ?? ""} onChange={(e) => setTime("ended_at", e.target.value)} />
      </div>

      {/* 場所 */}
      <div className="flex items-center gap-2">
        <label className="text-sm font-medium text-gray-700">場所</label>
        <input type="text" value={value.place ?? ""} onChange={(e) => setPlace(e.target.value)} placeholder="ジム名など" className="w-48 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" />
      </div>

      {/* 部位（プルダウン） */}
      <div className="flex items-center gap-2">
        <label className="text-sm font-medium text-gray-700">部位</label>
        <select value={value.workout_part?.id?.toString() ?? ""} onChange={(e) => setPartById(e.target.value)} className="w-40 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500 bg-white">
          <option value=""></option>
          {availableParts.map((part) => (
            <option key={part.id} value={part.id.toString()}>
              {part.name}
            </option>
          ))}
        </select>
      </div>

      {/* コンディション（1〜5） */}
      <div className="flex items-center gap-2">
        <label className="text-sm font-medium text-gray-700">コンディション</label>
        <div className="flex gap-1">
          {[1, 2, 3, 4, 5].map((n) => (
            <button key={n} type="button" onClick={() => setCondition(n as 1 | 2 | 3 | 4 | 5)} className={`w-10 h-10 rounded border text-sm font-medium transition-colors ${value.condition_level === n ? "bg-booking-600 text-white border-booking-600" : "border-gray-300 text-gray-700 hover:bg-gray-100"}`}>
              {n}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
};

export default WorkoutSessionMetaEditor;
