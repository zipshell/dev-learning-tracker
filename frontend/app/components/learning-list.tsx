import type { LearningPoint } from "../entries/queries";

type LearningListProps = {
  learningPoints: LearningPoint[];
  filterTag: string;
  filterOptions: string[];
  onFilterChange: (value: string) => void;
};

export function LearningList({
  learningPoints,
  filterTag,
  filterOptions,
  onFilterChange,
}: LearningListProps) {
  return (
    <section className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <h2 className="text-2xl font-semibold text-slate-900">Learned list</h2>
        <div className="flex flex-wrap items-center gap-3">
          <label htmlFor="filter-select" className="text-sm font-medium text-slate-700">
            Filter by tag:
          </label>
          <select
            id="filter-select"
            value={filterTag}
            onChange={(event) => onFilterChange(event.target.value)}
            className="rounded-2xl border border-slate-300 bg-white px-4 py-3 text-sm text-slate-900 outline-none transition focus:border-blue-500 focus:ring-4 focus:ring-blue-100"
          >
            <option value="">All tags</option>
            {filterOptions.map((tag) => (
              <option key={tag} value={tag}>
                {tag}
              </option>
            ))}
          </select>
        </div>
      </div>

      <div className="grid gap-6 sm:grid-cols-2 xl:grid-cols-3">
        {learningPoints.map((point, index) => (
          <article
            key={`${point.topic}-${index}`}
            className="rounded-[28px] border border-slate-200 bg-white/95 p-6 shadow-[0_18px_35px_rgba(15,23,42,0.08)]"
          >
            <h3 className="text-lg font-semibold text-slate-900">{point.topic}</h3>
            <div className="mt-4 flex flex-wrap gap-2">
              {point.tags.map((tag) => (
                <span
                  key={tag}
                  className="inline-flex rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-700"
                >
                  {tag}
                </span>
              ))}
            </div>
            <p className="mt-4 text-slate-600">{point.note}</p>
          </article>
        ))}
      </div>
    </section>
  );
}
