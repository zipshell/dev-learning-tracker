import { getCurrentUser, getUserStorageKey } from "./auth";

export type LearningPoint = {
  topic: string;
  tags: string[];
  note: string;
};

export function getSessionUser() {
  return getCurrentUser();
}

export function getStoredLearningPoints(email: string): LearningPoint[] {
  if (typeof window === "undefined") {
    return [];
  }

  const storageKey = getUserStorageKey(email);
  const saved = window.localStorage.getItem(storageKey);

  if (!saved) {
    return [];
  }

  try {
    return JSON.parse(saved) as LearningPoint[];
  } catch {
    return [];
  }
}

export function getFilterOptions(points: LearningPoint[]) {
  return Array.from(new Set(points.flatMap((point) => point.tags)));
}
