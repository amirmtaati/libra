"use client";
import { useState, useEffect } from "react";
import SearchBar from "../components/SearchBar";
import Tabs from "@/components/Tabs";
import BookList from "@/components/BookList";
import ShelfList from "@/components/ShelfList";

interface Book {
  id: number;
  title: string;
  author: string;
  // Add other book properties as needed
}

interface Shelf {
  id: number;
  name: string;
}

export default function HomePage() {
  const [search, setSearch] = useState("");
  const [selectedShelf, setSelectedShelf] = useState<number | null>(null);
  const [books, setBooks] = useState<Book[]>([]);
  const [shelves, setShelves] = useState<Shelf[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

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
	//        const shelvesResponse = await fetch('http://localhost:8080/api/shelves');
        //const shelvesData = await shelvesResponse.json();
        
        //if (!shelvesData.success) {
        //  throw new Error(shelvesData.error || 'Failed to fetch shelves');
	// }
        
        //setShelves(shelvesData.data || []);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error');
        console.error("Fetch error:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []); // Empty dependency array ensures this runs only once

  if (loading) {
    return <div className="p-8">Loading books...</div>;
  }

  if (error) {
    return <div className="p-8 text-red-500">Error: {error}</div>;
  }

  console.log(books)
  const filteredBooks = books.filter((book) =>
    book.Title.toLowerCase().includes(search.toLowerCase())
  );


  return (
    <div className="p-8">
      <SearchBar value={search} onChange={setSearch} />

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
            <ShelfList shelves={shelves} onSelectShelf={setSelectedShelf} />
          ),
        }}
      />
    </div>
  );
}
