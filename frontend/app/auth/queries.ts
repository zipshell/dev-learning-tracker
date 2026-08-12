import { getCurrentUser } from "../lib/auth";

export function getSessionUser() {
  return getCurrentUser();
}
