export default function BookList({ books }) {
  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      {books.map((book) => (
        <div key={book.ID} className="p-4 border rounded-lg shadow-sm">
          <h3 className="font-semibold">{book.Title}</h3>
          <p className="text-sm text-gray-500">{book.Author}</p>
        </div>
      ))}
    </div>
  );
}
