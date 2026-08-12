import Link from "next/link";

type DashboardHeaderProps = {
  userName: string;
};

export function DashboardHeader({ userName }: DashboardHeaderProps) {
  return (
    <header className="flex flex-col gap-4 rounded-[28px] border border-slate-200 bg-white/90 p-8 shadow-[0_24px_60px_rgba(15,23,42,0.08)] sm:flex-row sm:items-center sm:justify-between">
      <div>
        <p className="text-sm uppercase tracking-[0.35em] text-slate-400">Learning tracker</p>
        <h1 className="mt-3 text-4xl font-semibold tracking-tight text-slate-900 sm:text-5xl">
          Welcome back, {userName}
        </h1>
        <p className="mt-2 text-slate-600">
          Add new entries, filter your saved topics, and keep learning organized.
        </p>
      </div>

      <div className="flex flex-wrap items-center gap-3">
        <Link
          href="/signup"
          className="rounded-full border border-slate-300 bg-slate-50 px-5 py-3 text-sm font-semibold text-slate-800 transition hover:border-slate-400 hover:bg-slate-100"
        >
          Switch account
        </Link>
      </div>
    </header>
  );
}
