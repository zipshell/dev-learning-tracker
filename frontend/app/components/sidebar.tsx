import { SignOutButton } from "./sign-out-button";

export type FolderItem = {
  id: string;
  name: string;
};

type SidebarProps = {
  userName: string;
  userEmail: string;
};

const SAMPLE_FOLDERS = [
  { id: "all", name: "All notes" },
  { id: "frontend", name: "Frontend" },
  { id: "backend", name: "Backend" },
  { id: "design", name: "Design" },
  { id: "journal", name: "Journal" },
];

export function Sidebar({ userName, userEmail }: SidebarProps) {
  return (
    <aside className="w-full shrink-0 rounded-[32px] border border-slate-200 bg-white/90 p-6 shadow-[0_24px_60px_rgba(15,23,42,0.08)] lg:w-80">
      <div className="rounded-[24px] bg-slate-900 p-5 text-white">
        <p className="text-sm uppercase tracking-[0.35em] text-slate-400">Workspace</p>
        <h2 className="mt-3 text-xl font-semibold">{userName}</h2>
        <p className="mt-1 text-sm text-slate-300">{userEmail}</p>
      </div>

      <div className="mt-6">
        <div className="mb-3 flex items-center justify-between">
          <h3 className="text-sm font-semibold uppercase tracking-[0.24em] text-slate-500">
            Folders
          </h3>
          <button
            type="button"
            className="text-sm font-medium text-blue-600 transition hover:text-blue-700"
          >
            + New
          </button>
        </div>

        <nav className="space-y-2">
          {SAMPLE_FOLDERS.map((folder) => (
            <button
              key={folder.id}
              type="button"
              className="flex w-full items-center justify-between rounded-2xl px-3 py-3 text-left text-sm font-medium text-slate-700 transition hover:bg-slate-100"
            >
              <span>{folder.name}</span>
              <span className="text-slate-400">↗</span>
            </button>
          ))}
        </nav>
      </div>

      <div className="mt-8 border-t border-slate-200 pt-6">
        <SignOutButton />
      </div>
    </aside>
  );
}
