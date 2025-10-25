"use client";
import React, { useState } from "react";

export type ExerciseOption = { id: number; name: string };

type Props = {
  options: ExerciseOption[]; // 全件
  valueId: number | undefined;
  valueName: string;
  onPick: (picked: { id: number; name: string }) => void;
  onCreateNew: (name: string) => void;
};

const ExerciseSelect: React.FC<Props> = ({ options, valueId, valueName, onPick, onCreateNew }) => {
  const [mode, setMode] = useState<"select" | "custom">(!valueId && valueName ? "custom" : "select");
  const [custom, setCustom] = useState(valueName ?? "");

  return (
    <div className="flex items-center gap-2">
      {mode === "select" ? (
        <>
          <select
            className="w-full px-2 py-1 border border-gray-300 rounded"
            value={valueId ?? ""}
            onChange={(e) => {
              const id = Number(e.target.value);
              const opt = options.find((o) => o.id === id);
              if (opt) onPick({ id: opt.id, name: opt.name });
            }}
          >
            <option value="" disabled>
              種目を選択
            </option>
            {options.map((o) => (
              <option key={o.id} value={o.id}>
                {o.name}
              </option>
            ))}
          </select>
          <button type="button" className="px-2 py-1 text-xs border rounded" onClick={() => setMode("custom")}>
            カスタム
          </button>
        </>
      ) : (
        <>
          <input value={custom} onChange={(e) => setCustom(e.target.value)} placeholder="新規種目名" className="w-full px-2 py-1 border border-gray-300 rounded" />
          <button
            type="button"
            className="px-2 py-1 text-xs border rounded"
            onClick={() => {
              if (!custom.trim()) return;
              onCreateNew(custom.trim());
            }}
          >
            追加
          </button>
          <button type="button" className="px-2 py-1 text-xs border rounded" onClick={() => setMode("select")}>
            既存から選ぶ
          </button>
        </>
      )}
    </div>
  );
};

export default ExerciseSelect;
