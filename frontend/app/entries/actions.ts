import { getUserStorageKey } from "../lib/auth";
import type { LearningPoint } from "./queries";

export function persistLearningPoints(email: string, points: LearningPoint[]) {
  if (typeof window === "undefined") {
    return;
  }

  const storageKey = getUserStorageKey(email);
  window.localStorage.setItem(storageKey, JSON.stringify(points));
}

export function appendLearningPoint(
  email: string,
  point: LearningPoint,
  current: LearningPoint[],
) {
  const nextPoints = [...current, point];
  persistLearningPoints(email, nextPoints);
  return nextPoints;
}
