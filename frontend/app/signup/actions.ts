export async function signUp(
  name: string,
  email: string,
  password: string,
): Promise<{ success: boolean; error?: string }> {
  const normalized = email.trim();

  if (!email.trim()) {
    return { success: false, error: "Please provide an email address." };
  }

  if (!password) {
    return { success: false, error: "Please provide a password." };
  }

  const payload = {
    email: normalized,
    password,
  };

  try {
    const response = await fetch("/api/users", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      const errorText = await response.text();
      return {
        success: false,
        error: errorText || "Unable to create account.",
      };
    }

    return { success: true };
  } catch (error) {
    return {
      success: false,
      error: typeof error === "string" ? error : "Unable to create user",
    };
  }
}
