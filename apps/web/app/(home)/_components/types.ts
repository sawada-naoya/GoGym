export type Exercise = { id: string; name: string; isBuiltin?: boolean };

export type SessionMeta = {
  startedAt?: string | null;
  endedAt?: string | null;
  place?: string | null;
  muscle?: "" | "腕" | "胸" | "脚" | "肩" | "背中" | "体幹";
  condition?: 1 | 2 | 3 | 4 | 5 | null;
};
