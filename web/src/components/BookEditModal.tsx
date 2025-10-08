"use client";
import { useEffect, useState } from "react";

export interface EditableBookDto {
  ID: number;
  Title: string;
  Author: string;
  Tags?: string;
  Format?: string;
  FilePath?: string;
  Cover?: string;
}

export default function BookEditModal({
  book,
  open,
  onClose,
  onUpdated,
}: {
  book: EditableBookDto | null;
  open: boolean;
  onClose: () => void;
  onUpdated: (book: EditableBookDto) => void;
}) {
  const [title, setTitle] = useState("");
  const [author, setAuthor] = useState("");
  const [tags, setTags] = useState("");
  const [format, setFormat] = useState("");
  const [filepath, setFilepath] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [cover, setCover] = useState("");

  useEffect(() => {
    if (book) {
      setTitle(book.Title || "");
      setAuthor(book.Author || "");
      setTags(book.Tags || "");
      setFormat(book.Format || "");
      setFilepath(book.FilePath || "");
      setCover(book.Cover || "");
    }
  }, [book]);

  if (!open || !book) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      const form = new FormData();
      form.append("title", title);
      form.append("author", author);
      form.append("tags", tags);
      form.append("format", format);
      form.append("filepath", filepath);
      form.append("cover", cover);

      // POST to update endpoint: server has updateBook handler expecting form fields
      const res = await fetch(`http://localhost:8080/api/books/${book.ID}`, {
        method: "POST",
        body: form,
      });
      const contentType = res.headers.get("content-type") || "";
      const isJson = contentType.includes("application/json");
      const payload = isJson ? await res.json() : await res.text();

      if (!res.ok) {
        const msg = isJson ? (payload.error || JSON.stringify(payload)) : String(payload);
        throw new Error(msg || `Request failed with ${res.status}`);
      }
      const data = isJson ? payload : null;
      if (data && data.success === false) {
        throw new Error(data.error || "Failed to update book");
      }
      const updated = (data ? data.data : { ...book, Title: title, Author: author, Tags: tags, Format: format, FilePath: filepath, Cover: cover }) as EditableBookDto;
      onUpdated(updated);
      onClose();
    } catch (e: any) {
      setError(e?.message || "Failed to update book");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div className="w-full max-w-lg rounded-lg bg-white p-6 shadow-xl">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-lg font-semibold">Edit Book</h2>
          <button onClick={onClose} className="text-gray-500 hover:text-gray-700">✕</button>
        </div>
        {error && <div className="mb-3 rounded border border-red-300 bg-red-50 p-2 text-sm text-red-700">{error}</div>}
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="mb-1 block text-sm font-medium">Title</label>
            <input value={title} onChange={(e) => setTitle(e.target.value)} className="w-full rounded border p-2" required />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium">Author</label>
            <input value={author} onChange={(e) => setAuthor(e.target.value)} className="w-full rounded border p-2" required />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium">Tags</label>
            <input value={tags} onChange={(e) => setTags(e.target.value)} className="w-full rounded border p-2" placeholder="tag1, tag2" />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium">Format</label>
            <input value={format} onChange={(e) => setFormat(e.target.value)} className="w-full rounded border p-2 uppercase" placeholder="PDF, EPUB…" />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium">File path</label>
            <input value={filepath} onChange={(e) => setFilepath(e.target.value)} className="w-full rounded border p-2" placeholder="/path/to/file" />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium">Cover URL (optional)</label>
            <input value={cover} onChange={(e) => setCover(e.target.value)} className="w-full rounded border p-2" placeholder="https://..." />
          </div>
          <div className="flex justify-end gap-2 pt-2">
            <button type="button" onClick={onClose} className="rounded border px-4 py-2">Cancel</button>
            <button disabled={submitting} type="submit" className="rounded bg-blue-600 px-4 py-2 text-white disabled:opacity-50">
              {submitting ? "Saving…" : "Save"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}


