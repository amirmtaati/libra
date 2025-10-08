"use client";
import { useState, useEffect } from "react";

export default function RescanModal({
  open,
  onClose,
  onTriggered,
}: {
  open: boolean;
  onClose: () => void;
  onTriggered: () => void;
}) {
  const [path, setPath] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!open) {
      setPath("");
      setError(null);
      setSubmitting(false);
    }
  }, [open]);

  if (!open) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      const res = await fetch("http://localhost:8080/api/scan", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ path }),
      });
      const contentType = res.headers.get("content-type") || "";
      const isJson = contentType.includes("application/json");
      const payload = isJson ? await res.json() : await res.text();
      if (!res.ok) {
        const msg = isJson ? (payload.error || JSON.stringify(payload)) : String(payload);
        throw new Error(msg || `Request failed with ${res.status}`);
      }
      onTriggered();
      onClose();
    } catch (e: any) {
      setError(e?.message || "Failed to trigger scan");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div className="w-full max-w-lg rounded-lg bg-white p-6 shadow-xl">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-lg font-semibold">Re-scan Library</h2>
          <button onClick={onClose} className="text-gray-500 hover:text-gray-700">✕</button>
        </div>
        {error && <div className="mb-3 rounded border border-red-300 bg-red-50 p-2 text-sm text-red-700">{error}</div>}
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="mb-1 block text-sm font-medium">Directory path (optional)</label>
            <input value={path} onChange={(e) => setPath(e.target.value)} className="w-full rounded border p-2" placeholder="/path/to/dir (leave empty to reuse)" />
          </div>
          <div className="flex justify-end gap-2 pt-2">
            <button type="button" onClick={onClose} className="rounded border px-4 py-2">Cancel</button>
            <button disabled={submitting} type="submit" className="rounded bg-blue-600 px-4 py-2 text-white disabled:opacity-50">{submitting ? "Starting…" : "Start Scan"}</button>
          </div>
        </form>
      </div>
    </div>
  );
}


