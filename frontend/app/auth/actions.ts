import { logout as logoutUser } from "../lib/auth";

export function signOutUser() {
  logoutUser();
}
