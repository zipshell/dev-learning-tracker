"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { getSessionUser } from "../auth/queries";
import { login } from "./actions";

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isCheckingSession, setIsCheckingSession] = useState(true);

  const authGuard = function () {
    const user = getSessionUser();
    if (user) {
      router.replace("/");
      return;
    }

    setIsCheckingSession(false);
  };

  useEffect(authGuard, [router]);

  const handleSubmit = async (event: React.SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError("");

    const result = await login(email, password);
    if (!result.success) {
      setError(result.error ?? "Unable to sign in.");
      return;
    }

    router.push("/");
  };

  if (isCheckingSession) {
    return (
      <main className="min-h-screen bg-slate-50 px-4 py-16">
        <div className="mx-auto w-full max-w-xl rounded-[28px] border border-slate-200 bg-white p-10 text-center shadow-[0_24px_60px_rgba(15,23,42,0.08)]">
          <p className="text-slate-500">Checking session…</p>
        </div>
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-slate-50 px-4 py-16">
      <div className="mx-auto w-full max-w-xl rounded-[32px] border border-slate-200 bg-white p-10 shadow-[0_24px_60px_rgba(15,23,42,0.08)]">
        <div className="mb-8 text-center">
          <p className="text-sm uppercase tracking-[0.35em] text-slate-400">
            Learning tracker
          </p>
          <h1 className="mt-4 text-3xl font-semibold text-slate-900">
            Welcome back
          </h1>
          <p className="mt-2 text-slate-600">
            Sign in to manage your learning notes and saved entries.
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid gap-2">
            <label
              htmlFor="email"
              className="text-sm font-medium text-slate-700"
            >
              Email
            </label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              className="w-full rounded-2xl border border-slate-300 bg-slate-50 px-4 py-3 text-slate-900 outline-none transition focus:border-blue-500 focus:ring-4 focus:ring-blue-100"
              placeholder="you@example.com"
              autoComplete="email"
            />
          </div>

          <div className="grid gap-2">
            <label
              htmlFor="password"
              className="text-sm font-medium text-slate-700"
            >
              Password
            </label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              className="w-full rounded-2xl border border-slate-300 bg-slate-50 px-4 py-3 text-slate-900 outline-none transition focus:border-blue-500 focus:ring-4 focus:ring-blue-100"
              placeholder="Enter your password"
              autoComplete="current-password"
            />
          </div>

          {error ? (
            <div className="rounded-2xl bg-rose-100 px-4 py-3 text-sm text-rose-700">
              {error}
            </div>
          ) : null}

          <button
            type="submit"
            className="w-full rounded-full bg-gradient-to-r from-blue-600 to-teal-500 px-6 py-3 text-sm font-semibold text-white shadow-lg shadow-teal-200/30 transition hover:-translate-y-0.5"
          >
            Sign in
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-slate-600">
          New here?{" "}
          <Link
            href="/signup"
            className="font-semibold text-slate-900 hover:text-blue-600"
          >
            Create an account
          </Link>
        </p>
      </div>
    </main>
  );
}
