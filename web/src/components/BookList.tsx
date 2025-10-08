export interface BookDto {
  ID: number;
  Title: string;
  Author: string;
  Tags?: string;
  Format?: string;
  FilePath?: string;
  Cover?: string;
}

export default function BookList({ books, onEdit, onDelete }: { books: BookDto[]; onEdit?: (book: BookDto) => void; onDelete?: (book: BookDto) => void; }) {
  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      {books.map((book) => (
        <div key={book.ID} className="border rounded-lg shadow-sm flex flex-col">
          {book.Cover ? (
            <div className="w-full rounded-t-lg bg-white overflow-hidden h-56 md:h-64 flex items-center justify-center">
              <img src={book.Cover} alt={book.Title} className="max-h-full max-w-full object-contain" />
            </div>
          ) : (
            <div className="w-full rounded-t-lg bg-gray-100 h-56 md:h-64 flex items-center justify-center text-gray-400 text-sm">No cover</div>
          )}
          <div className="p-4 flex flex-col gap-2">
            <h3 className="font-semibold">{book.Title}</h3>
            <p className="text-sm text-gray-500">{book.Author}</p>
            <div className="mt-1 flex gap-3">
              {onEdit && (
                <button onClick={() => onEdit(book)} className="self-start text-sm text-blue-600 hover:underline">Edit</button>
              )}
              {onDelete && (
                <button onClick={() => onDelete(book)} className="self-start text-sm text-red-600 hover:underline">Delete</button>
              )}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
