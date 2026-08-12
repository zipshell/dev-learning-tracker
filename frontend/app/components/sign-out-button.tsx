"use client";

import { useRouter } from "next/navigation";
import { signOutUser } from "../auth/actions";

export function SignOutButton() {
  const router = useRouter();

  const handleSignOut = () => {
    signOutUser();
    router.push("/login");
  };

  return (
    <button
      type="button"
      onClick={handleSignOut}
      className="w-full rounded-full bg-gradient-to-r from-rose-500 to-pink-500 px-4 py-3 text-sm font-semibold text-white shadow-lg shadow-pink-200/30 transition hover:-translate-y-0.5"
    >
      Sign out
    </button>
  );
}
