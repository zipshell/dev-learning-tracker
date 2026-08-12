import type { FormEvent, KeyboardEvent } from "react";

type LearningFormProps = {
  topic: string;
  tagInput: string;
  selectedTags: string[];
  note: string;
  suggestions: string[];
  onTopicChange: (value: string) => void;
  onTagInputChange: (value: string) => void;
  onTagKeyDown: (event: KeyboardEvent<HTMLInputElement>) => void;
  onTagRemove: (index: number) => void;
  onTagSelect: (tag: string) => void;
  onNoteChange: (value: string) => void;
  onSubmit: (event: FormEvent<HTMLFormElement>) => void;
};

export function LearningForm({
  topic,
  tagInput,
  selectedTags,
  note,
  suggestions,
  onTopicChange,
  onTagInputChange,
  onTagKeyDown,
  onTagRemove,
  onTagSelect,
  onNoteChange,
  onSubmit,
}: LearningFormProps) {
  return (
    <form
      onSubmit={onSubmit}
      className="grid gap-6 rounded-[28px] border border-slate-200 bg-white/90 p-8 shadow-[0_24px_60px_rgba(15,23,42,0.08)]"
    >
      <h2 className="text-center text-xl font-semibold text-slate-900">
        Add new learning points
      </h2>

      <div className="grid gap-2">
        <label htmlFor="learning-point" className="text-sm font-medium text-slate-700">
          Learning point
        </label>
        <input
          id="learning-point"
          value={topic}
          onChange={(event) => onTopicChange(event.target.value)}
          className="w-full rounded-2xl border border-slate-300 bg-slate-50 px-4 py-3 text-slate-900 outline-none transition focus:border-blue-500 focus:ring-4 focus:ring-blue-100"
        />
      </div>

      <div className="grid gap-2">
        <label htmlFor="new-tag" className="text-sm font-medium text-slate-700">
          Tag
        </label>
        <div className="flex flex-wrap items-center gap-2 rounded-2xl border border-slate-300 bg-slate-50 p-3">
          <input
            id="new-tag"
            value={tagInput}
            autoComplete="off"
            placeholder="Type a tag and press Enter"
            onChange={(event) => onTagInputChange(event.target.value)}
            onKeyDown={onTagKeyDown}
            className="flex-1 min-w-[140px] bg-transparent text-slate-900 outline-none placeholder:text-slate-400"
          />
          <div className="flex flex-wrap gap-2" aria-live="polite">
            {selectedTags.map((tag, index) => (
              <span
                key={`${tag}-${index}`}
                className="inline-flex items-center gap-2 rounded-full bg-blue-100 px-3 py-1 text-sm font-medium text-blue-700"
              >
                {tag}
                <button
                  type="button"
                  onClick={() => onTagRemove(index)}
                  aria-label={`Remove ${tag}`}
                  className="text-blue-700 transition hover:text-blue-900"
                >
                  ×
                </button>
              </span>
            ))}
          </div>
        </div>

        <ul
          className={`mt-2 max-h-48 overflow-y-auto rounded-2xl border border-slate-200 bg-white shadow-lg ${suggestions.length === 0 ? "hidden" : ""}`}
        >
          {suggestions.map((tag) => (
            <li
              key={tag}
              className="cursor-pointer px-4 py-3 text-slate-700 hover:bg-slate-100"
              onClick={() => onTagSelect(tag)}
            >
              {tag}
            </li>
          ))}
        </ul>
      </div>

      <div className="grid gap-2">
        <label htmlFor="new-note" className="text-sm font-medium text-slate-700">
          Note
        </label>
        <textarea
          id="new-note"
          rows={5}
          value={note}
          onChange={(event) => onNoteChange(event.target.value)}
          className="w-full rounded-2xl border border-slate-300 bg-slate-50 px-4 py-3 text-slate-900 outline-none transition focus:border-blue-500 focus:ring-4 focus:ring-blue-100"
        />
      </div>

      <button
        type="submit"
        className="mx-auto rounded-full bg-gradient-to-r from-blue-600 to-teal-500 px-8 py-3 text-sm font-semibold text-white shadow-lg shadow-teal-200/30 transition hover:-translate-y-0.5"
      >
        Add
      </button>
    </form>
  );
}
