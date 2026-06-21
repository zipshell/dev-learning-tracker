"use client";

import { useEffect, useMemo, useRef, useState } from "react";

const AVAILABLE_TAGS = [
  "frontend",
  "javascript",
  "react",
  "css",
  "html",
  "design",
  "backend",
  "accessibility",
];

const STORAGE_KEY = "learningPoints";

type LearningPoint = {
  topic: string;
  tags: string[];
  note: string;
};

export default function Home() {
  const [topic, setTopic] = useState("");
  const [tagInput, setTagInput] = useState("");
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [note, setNote] = useState("");
  const [learningPoints, setLearningPoints] = useState<LearningPoint[]>([]);
  const [filterTag, setFilterTag] = useState("");
  const hasLoadedStorage = useRef(false);

  useEffect(() => {
    const saved = localStorage.getItem(STORAGE_KEY);
    if (saved) {
      try {
        setLearningPoints(JSON.parse(saved));
      } catch (error) {
        console.warn("Failed to parse saved learning points", error);
      }
    }
    hasLoadedStorage.current = true;
  }, []);

  useEffect(() => {
    if (!hasLoadedStorage.current) return;
    localStorage.setItem(STORAGE_KEY, JSON.stringify(learningPoints));
  }, [learningPoints]);

  const filteredLearningPoints = useMemo(
    () =>
      filterTag
        ? learningPoints.filter((point) => point.tags.includes(filterTag))
        : learningPoints,
    [filterTag, learningPoints],
  );

  const filterOptions = useMemo(
    () => Array.from(new Set(learningPoints.flatMap((point) => point.tags))),
    [learningPoints],
  );

  const suggestions = useMemo(() => {
    const query = tagInput.trim().toLowerCase();
    return AVAILABLE_TAGS.filter(
      (tag) => tag.includes(query) && !selectedTags.includes(tag),
    );
  }, [selectedTags, tagInput]);

  const addTag = (value: string) => {
    const tag = value.trim().replace(/,+$/, "");
    if (!tag || selectedTags.includes(tag)) return;

    setSelectedTags((current) => [...current, tag]);
    setTagInput("");
  };

  const removeTag = (index: number) => {
    setSelectedTags((current) => current.filter((_, idx) => idx !== index));
  };

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (!topic.trim() || !note.trim() || selectedTags.length === 0) {
      return;
    }

    setLearningPoints((current) => [
      ...current,
      {
        topic: topic.trim(),
        tags: selectedTags,
        note: note.trim(),
      },
    ]);

    setTopic("");
    setTagInput("");
    setSelectedTags([]);
    setNote("");
  };

  const handleTagKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Enter" || event.key === ",") {
      event.preventDefault();
      addTag(tagInput);
    } else if (event.key === "Backspace" && !tagInput && selectedTags.length) {
      removeTag(selectedTags.length - 1);
    }
  };

  return (
    <main className="min-h-screen bg-slate-50 py-10 px-4 sm:px-6 lg:px-8">
      <div className="mx-auto w-full max-w-6xl">
        <h1 className="text-4xl sm:text-5xl font-semibold tracking-tight text-slate-900 mb-8">
          The Learning Tracker
        </h1>

        <form
          onSubmit={handleSubmit}
          className="grid gap-6 rounded-[28px] border border-slate-200 bg-white/90 p-8 shadow-[0_24px_60px_rgba(15,23,42,0.08)]"
        >
          <h2 className="text-xl font-semibold text-slate-900 text-center">
            Add new learning points
          </h2>

          <div className="grid gap-2">
            <label htmlFor="learning-point" className="text-sm font-medium text-slate-700">
              Learning point
            </label>
            <input
              id="learning-point"
              value={topic}
              onChange={(event) => setTopic(event.target.value)}
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
                onChange={(event) => setTagInput(event.target.value)}
                onKeyDown={handleTagKeyDown}
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
                      onClick={() => removeTag(index)}
                      aria-label={`Remove ${tag}`}
                      className="text-blue-700 transition hover:text-blue-900"
                    >
                      &times;
                    </button>
                  </span>
                ))}
              </div>
            </div>

            <ul
              className={`mt-2 max-h-48 overflow-y-auto rounded-2xl border border-slate-200 bg-white shadow-lg ${
                suggestions.length === 0 ? "hidden" : ""
              }`}
            >
              {suggestions.map((tag) => (
                <li
                  key={tag}
                  className="cursor-pointer px-4 py-3 text-slate-700 hover:bg-slate-100"
                  onClick={() => addTag(tag)}
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
              onChange={(event) => setNote(event.target.value)}
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

        <section className="mt-10 space-y-6">
          <div className="flex flex-wrap items-center justify-between gap-4">
            <h2 className="text-2xl font-semibold text-slate-900">Learned list</h2>
            <div className="flex flex-wrap items-center gap-3">
              <label htmlFor="filter-select" className="text-sm font-medium text-slate-700">
                Filter by tag:
              </label>
              <select
                id="filter-select"
                value={filterTag}
                onChange={(event) => setFilterTag(event.target.value)}
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
            {filteredLearningPoints.map((point, index) => (
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
      </div>
    </main>
  );
}
