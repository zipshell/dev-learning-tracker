import React from "react";

export const getSessionUser = React.cache(async function () {
  try {
    const response = await fetch("/api/users/me", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    console.log("response", response);

    if (!response.ok) {
      const errorText = await response.text();
      return {
        success: false,
        error: errorText || "Invalid session",
      };
    }

    const data = await response.json();

    return data;
  } catch (error) {
    console.log("error", error);
    return {
      success: false,
      error: typeof error === "string" ? error : "Invalid session",
    };
  }
});
