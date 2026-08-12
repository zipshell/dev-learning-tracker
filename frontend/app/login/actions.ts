export async function login(
  email: string,
  password: string,
): Promise<{ success: boolean; error?: string }> {
  if (typeof window === "undefined") {
    return { success: false, error: "Unable to access browser storage." };
  }

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
    const response = await fetch("/api/auth/login", {
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
        error: errorText || "Unable to log in",
      };
    }

    const data = await response.json();
    console.log(data);

    return { success: true };
  } catch (error) {
    return {
      success: false,
      error: typeof error === "string" ? error : "Unable to log in",
    };
  }
}
