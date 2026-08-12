"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { LearningForm } from "./learning-form";
import { LearningList } from "./learning-list";
import { appendLearningPoint, persistLearningPoints } from "../entries/actions";
import {
  getFilterOptions,
  type LearningPoint,
} from "../entries/queries";

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

type LearningWorkspaceProps = {
  userEmail: string;
};

export function LearningWorkspace({ userEmail }: LearningWorkspaceProps) {
  const [topic, setTopic] = useState("");
  const [tagInput, setTagInput] = useState("");
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [note, setNote] = useState("");
  const [learningPoints, setLearningPoints] = useState<LearningPoint[]>([]);
  const [filterTag, setFilterTag] = useState("");
  const hasLoadedStorage = useRef(false);

  useEffect(() => {
    if (!hasLoadedStorage.current) return;
    persistLearningPoints(userEmail, learningPoints);
  }, [learningPoints, userEmail]);

  const filteredLearningPoints = useMemo(
    () =>
      filterTag
        ? learningPoints.filter((point) => point.tags.includes(filterTag))
        : learningPoints,
    [filterTag, learningPoints],
  );

  const filterOptions = useMemo(
    () => getFilterOptions(learningPoints),
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

    const nextPoint: LearningPoint = {
      topic: topic.trim(),
      tags: selectedTags,
      note: note.trim(),
    };

    setLearningPoints((current) =>
      appendLearningPoint(userEmail, nextPoint, current),
    );

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
    <>
      <LearningForm
        topic={topic}
        tagInput={tagInput}
        selectedTags={selectedTags}
        note={note}
        suggestions={suggestions}
        onTopicChange={setTopic}
        onTagInputChange={setTagInput}
        onTagKeyDown={handleTagKeyDown}
        onTagRemove={removeTag}
        onTagSelect={addTag}
        onNoteChange={setNote}
        onSubmit={handleSubmit}
      />

      <LearningList
        learningPoints={filteredLearningPoints}
        filterTag={filterTag}
        filterOptions={filterOptions}
        onFilterChange={setFilterTag}
      />
    </>
  );
}
