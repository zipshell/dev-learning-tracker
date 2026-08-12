import { login as loginUser, logout as logoutUser } from "./auth";

export type LearningPoint = {
  topic: string;
  tags: string[];
  note: string;
};

export function submitLogin(email: string, password: string) {
  return loginUser(email, password);
}

export function signOutUser() {
  logoutUser();
}
