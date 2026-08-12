export type StoredUser = {
  name: string;
  email: string;
  password: string;
};

export type AuthUser = {
  name: string;
  email: string;
};

const USERS_KEY = "learningTrackerUsers";
const SESSION_KEY = "learningTrackerSession";

function normalizeEmail(email: string) {
  return email.trim().toLowerCase();
}

function getStoredUsers(): Record<string, StoredUser> {
  if (typeof window === "undefined") {
    return {};
  }

  const raw = localStorage.getItem(USERS_KEY);
  if (!raw) return {};

  try {
    return JSON.parse(raw) as Record<string, StoredUser>;
  } catch {
    return {};
  }
}

export function login(
  email: string,
  password: string,
): { success: boolean; error?: string } {
  if (typeof window === "undefined") {
    return { success: false, error: "Unable to access browser storage." };
  }

  const normalized = normalizeEmail(email);
  const users = getStoredUsers();
  const storedUser = users[normalized];

  if (!storedUser) {
    return { success: false, error: "No account found with this email." };
  }

  if (storedUser.password !== password) {
    return { success: false, error: "Invalid email or password." };
  }

  localStorage.setItem(SESSION_KEY, storedUser.email);
  return { success: true };
}

export function logout() {
  if (typeof window === "undefined") return;
  localStorage.removeItem(SESSION_KEY);
}
