"use client";
import { useState, useEffect } from "react";
import SearchBar from "../components/SearchBar";
import Tabs from "../components/Tabs";
import BookList, { BookDto } from "../components/BookList";
import ShelfList from "../components/ShelfList";

interface Shelf {
  id: number;
  name: string;
}

export default function HomePage() {
  const [search, setSearch] = useState("");
  const [selectedShelf, setSelectedShelf] = useState<number | null>(null);
  const [books, setBooks] = useState<BookDto[]>([]);
  const [shelves, setShelves] = useState<Shelf[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [newShelfName, setNewShelfName] = useState("");
  const [creatingShelf, setCreatingShelf] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Fetch books
        const booksResponse = await fetch('http://localhost:8080/api/books');
        const booksData = await booksResponse.json();
        
        if (!booksData.success) {
          throw new Error(booksData.error || 'Failed to fetch books');
        }
        
        setBooks(booksData.data || []);

        // Fetch shelves
        const shelvesResponse = await fetch('http://localhost:8080/api/shelves');
        const shelvesData = await shelvesResponse.json();
        if (!shelvesData.success) {
          throw new Error(shelvesData.error || 'Failed to fetch shelves');
        }
        const mappedShelves: Shelf[] = (shelvesData.data || []).map((s: any) => ({ id: s.ID, name: s.Name }));
        setShelves(mappedShelves);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error');
        console.error("Fetch error:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []); // Empty dependency array ensures this runs only once

  const handleCreateShelf = async () => {
    const name = newShelfName.trim();
    if (!name) return;
    try {
      setCreatingShelf(true);
      const res = await fetch('http://localhost:8080/api/shelves', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name }),
      });
      const data = await res.json();
      if (!data.success) {
        throw new Error(data.error || 'Failed to create shelf');
      }
      const created = data.data;
      const newShelf: Shelf = { id: created.ID, name: created.Name };
      setShelves((prev) => [...prev, newShelf]);
      setNewShelfName("");
    } catch (e) {
      console.error('Create shelf error:', e);
      alert(e instanceof Error ? e.message : 'Failed to create shelf');
    } finally {
      setCreatingShelf(false);
    }
  };

  if (loading) {
    return <div className="p-8">Loading books...</div>;
  }

  if (error) {
    return <div className="p-8 text-red-500">Error: {error}</div>;
  }

  console.log(books)
  const filteredBooks = books.filter((book) =>
    // @ts-ignore backend returns Title
    book.Title.toLowerCase().includes(search.toLowerCase())
  );


  return (
    <div className="p-8">
      <SearchBar value={search} onChange={setSearch} onClear={() => setSearch("")} />

      <Tabs
        tabs={{
          "All Books": <BookList books={filteredBooks} />,
          Shelves: selectedShelf ? (
            <div>
              <button
                className="mb-4 text-blue-500 underline"
                onClick={() => setSelectedShelf(null)}
              >
                ← Back to shelves
              </button>
              <BookList books={books} /> {/* Will update later with actual shelf books */}
            </div>
          ) : (
            <div className="space-y-4">
              <div className="flex items-center gap-2">
                <input
                  type="text"
                  value={newShelfName}
                  onChange={(e) => setNewShelfName(e.target.value)}
                  placeholder="New shelf name"
                  className="flex-1 p-2 border rounded"
                />
                <button
                  onClick={handleCreateShelf}
                  disabled={creatingShelf || newShelfName.trim().length === 0}
                  className="px-3 py-2 bg-blue-600 text-white rounded disabled:opacity-50"
                >
                  {creatingShelf ? 'Creating…' : 'Add shelf'}
                </button>
              </div>
              <ShelfList shelves={shelves} onSelectShelf={setSelectedShelf} />
            </div>
          ),
        }}
      />
    </div>
  );
}
