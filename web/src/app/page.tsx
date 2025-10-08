"use client";
import { useState, useEffect } from "react";
import SearchBar from "../components/SearchBar";
import Tabs from "../components/Tabs";
import BookList, { BookDto } from "../components/BookList";
import BookEditModal, { EditableBookDto } from "../components/BookEditModal";
import ConfirmDialog from "../components/ConfirmDialog";
import RescanModal from "../components/RescanModal";
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
  const [shelfBooks, setShelfBooks] = useState<BookDto[]>([]);
  const [loadingShelfBooks, setLoadingShelfBooks] = useState(false);
  const [shelfBooksError, setShelfBooksError] = useState<string | null>(null);
  const [editingBook, setEditingBook] = useState<EditableBookDto | null>(null);
  const [editOpen, setEditOpen] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [bookPendingDelete, setBookPendingDelete] = useState<BookDto | null>(null);
  const [rescanOpen, setRescanOpen] = useState(false);
  const [rescanPathModalOpen, setRescanPathModalOpen] = useState(false);

  const fetchBooks = async () => {
    const booksResponse = await fetch('http://localhost:8080/api/books');
    const booksData = await booksResponse.json();
    if (!booksData.success) throw new Error(booksData.error || 'Failed to fetch books');
    setBooks(booksData.data || []);
  };

  const fetchShelves = async () => {
    const shelvesResponse = await fetch('http://localhost:8080/api/shelves');
    const shelvesData = await shelvesResponse.json();
    if (!shelvesData.success) {
      const errMsg = (shelvesData.error || '').toLowerCase();
      if (errMsg.includes('record not found')) {
        setShelves([]);
        return;
      }
      throw new Error(shelvesData.error || 'Failed to fetch shelves');
    }
    const mappedShelves: Shelf[] = (shelvesData.data || []).map((s: any) => ({ id: s.ID, name: s.Name }));
    setShelves(mappedShelves);
  };

  const fetchAll = async () => {
    try {
      await Promise.all([fetchBooks(), fetchShelves()]);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
      console.error('Fetch error:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAll();
    const interval = setInterval(fetchAll, 10000);
    const onFocus = () => fetchAll();
    window.addEventListener('focus', onFocus);
    return () => {
      clearInterval(interval);
      window.removeEventListener('focus', onFocus);
    };
  }, []);

  const fetchShelfBooks = async () => {
    if (selectedShelf === null) return;
    try {
      setLoadingShelfBooks(true);
      setShelfBooksError(null);
      const res = await fetch(`http://localhost:8080/api/shelves/${selectedShelf}/books`);
      const data = await res.json();
      if (!data.success) {
        throw new Error(data.error || 'Failed to fetch shelf books');
      }
      setShelfBooks(data.data || []);
    } catch (e) {
      setShelfBooksError(e instanceof Error ? e.message : 'Unknown error');
      setShelfBooks([]);
    } finally {
      setLoadingShelfBooks(false);
    }
  };

  useEffect(() => {
    if (selectedShelf === null) return;
    fetchShelfBooks();
    const interval = setInterval(fetchShelfBooks, 10000);
    const onFocus = () => fetchShelfBooks();
    window.addEventListener('focus', onFocus);
    return () => {
      clearInterval(interval);
      window.removeEventListener('focus', onFocus);
    };
  }, [selectedShelf]);

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

  const handleAddBookToShelf = async (bookId: number) => {
    if (selectedShelf === null) return;
    try {
      const res = await fetch(`http://localhost:8080/api/shelves/${selectedShelf}/books/${bookId}`, { method: 'POST' });
      const data = await res.json();
      if (!data.success) {
        throw new Error(data.error || 'Failed to add book to shelf');
      }
      const book = books.find((b) => b.ID === bookId);
      if (book) {
        setShelfBooks((prev) => [...prev, book]);
      }
    } catch (e) {
      console.error('Add book error:', e);
      alert(e instanceof Error ? e.message : 'Failed to add book');
    }
  };

  const handleRemoveBookFromShelf = async (bookId: number) => {
    if (selectedShelf === null) return;
    try {
      const res = await fetch(`http://localhost:8080/api/shelves/${selectedShelf}/books/${bookId}`, { method: 'DELETE' });
      const data = await res.json();
      if (!data.success) {
        throw new Error(data.error || 'Failed to remove book from shelf');
      }
      setShelfBooks((prev) => prev.filter((b) => b.ID !== bookId));
    } catch (e) {
      console.error('Remove book error:', e);
      alert(e instanceof Error ? e.message : 'Failed to remove book');
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

  const currentShelfName = selectedShelf ? shelves.find((s) => s.id === selectedShelf)?.name : undefined;
  const booksNotInShelf = selectedShelf !== null ? books.filter((b) => !shelfBooks.some((sb) => sb.ID === b.ID)) : [];

  return (
    <div className="p-8">
      <div className="mb-6 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div className="flex-1">
          <SearchBar value={search} onChange={setSearch} onClear={() => setSearch("")} />
        </div>
        <div className="flex-shrink-0">
          <button
            onClick={() => setRescanOpen(true)}
            className="px-4 py-2 rounded-lg bg-indigo-600 text-white"
          >
            Re-scan Library
          </button>
        </div>
      </div>

      <Tabs
        tabs={{
          "All Books": (
            <>
              <BookList
                books={filteredBooks}
                onEdit={(book) => {
                  setEditingBook(book as EditableBookDto);
                  setEditOpen(true);
                }}
                onDelete={(book) => {
                  setBookPendingDelete(book);
                  setConfirmOpen(true);
                }}
              />
              <BookEditModal
                book={editingBook}
                open={editOpen}
                onClose={() => setEditOpen(false)}
                onUpdated={(updated) => {
                  setBooks((prev) => prev.map((b) => (b.ID === updated.ID ? { ...b, ...updated } : b)));
                  setShelfBooks((prev) => prev.map((b) => (b.ID === updated.ID ? { ...b, ...updated } : b)));
                }}
              />
              <ConfirmDialog
                open={confirmOpen}
                title="Delete book?"
                message="This will permanently remove the book. Are you sure?"
                confirmText="Delete"
                onCancel={() => { setConfirmOpen(false); setBookPendingDelete(null); }}
                onConfirm={async () => {
                  if (!bookPendingDelete) return;
                  try {
                    const res = await fetch(`http://localhost:8080/api/books/${bookPendingDelete.ID}`, { method: 'DELETE' });
                    const data = await res.json();
                    if (!res.ok || !data.success) throw new Error(data.error || 'Failed to delete');
                    setBooks((prev) => prev.filter((b) => b.ID !== bookPendingDelete.ID));
                    setShelfBooks((prev) => prev.filter((b) => b.ID !== bookPendingDelete.ID));
                  } catch (e) {
                    alert(e instanceof Error ? e.message : 'Failed to delete');
                  } finally {
                    setConfirmOpen(false);
                    setBookPendingDelete(null);
                  }
                }}
              />
              <ConfirmDialog
                open={rescanOpen}
                title="Re-scan library?"
                message="Scanning can take time. Do you want to configure the directory before starting?"
                confirmText="Configure & Start"
                cancelText="Cancel"
                onCancel={() => setRescanOpen(false)}
                onConfirm={() => {
                  setRescanOpen(false);
                  setTimeout(() => {
                    setRescanPathModalOpen(true);
                  }, 0);
                }}
              />
              <RescanModal
                open={rescanPathModalOpen}
                onClose={() => setRescanPathModalOpen(false)}
                onTriggered={() => {
                  // Immediately refresh after scan starts
                  setTimeout(() => {
                    fetchAll();
                    if (selectedShelf !== null) fetchShelfBooks();
                  }, 1000); // Give scan a moment to start
                }}
              />
            </>
          ),
          Shelves: selectedShelf ? (
            <div className="space-y-4">
              <button
                className="mb-2 text-blue-500 underline"
                onClick={() => setSelectedShelf(null)}
              >
                ← Back to shelves
              </button>
              <h2 className="text-xl font-semibold">{currentShelfName || `Shelf #${selectedShelf}`}</h2>
              {loadingShelfBooks ? (
                <div>Loading shelf books...</div>
              ) : shelfBooksError ? (
                <div className="text-red-500">{shelfBooksError}</div>
              ) : (
                <div className="space-y-6">
                  {shelfBooks.length > 0 ? (
                    <div>
                      <h3 className="font-medium mb-2">Books in this shelf</h3>
                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                        {shelfBooks.map((book) => (
                          <div key={book.ID} className="p-4 border rounded-lg shadow-sm flex flex-col gap-1">
                            <div className="font-semibold">{book.Title}</div>
                            <div className="text-sm text-gray-500">{book.Author}</div>
                            <button
                              onClick={() => handleRemoveBookFromShelf(book.ID)}
                              className="mt-2 text-sm text-red-600 hover:underline self-start"
                            >
                              Remove
                            </button>
                          </div>
                        ))}
                      </div>
                    </div>
                  ) : (
                    <div className="text-gray-600">This shelf is empty.</div>
                  )}
                  <div>
                    <h3 className="font-medium mb-2">Add books to this shelf</h3>
                    {booksNotInShelf.length === 0 ? (
                      <div className="text-gray-600">All books are already in this shelf.</div>
                    ) : (
                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                        {booksNotInShelf.map((book) => (
                          <div key={book.ID} className="p-4 border rounded-lg shadow-sm flex flex-col gap-1">
                            <div className="font-semibold">{book.Title}</div>
                            <div className="text-sm text-gray-500">{book.Author}</div>
                            <button
                              onClick={() => handleAddBookToShelf(book.ID)}
                              className="mt-2 text-sm text-blue-600 hover:underline self-start"
                            >
                              Add
                            </button>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </div>
              )}
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
