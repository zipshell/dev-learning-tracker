import { Suspense } from "react";
import { DashboardHeader } from "./components/dashboard-header";
import { LearningWorkspace } from "./components/learning-workspace";
import { Sidebar } from "./components/sidebar";
import { getSessionUser } from "./utils/utils";

export default async function Home() {
  const user = await getSessionUser()

  return (
    <main className="min-h-screen bg-slate-50 py-10 px-4 sm:px-6 lg:px-8">
      <div className="mx-auto flex w-full max-w-7xl flex-col gap-6 lg:flex-row">
        <Suspense>
          <Sidebar
            userName={user.name}
            userEmail={user.email}
          />
        </Suspense>

        {/* <div className="flex-1 space-y-6">
          <DashboardHeader userName={user.name} />

          <LearningWorkspace userEmail={user.email} />
        </div> */}
      </div>
    </main>
  );
}
