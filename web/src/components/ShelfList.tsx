export default function ShelfList({ shelves, onSelectShelf }) {
  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      {shelves.map((shelf) => (
        <button
          key={shelf.id}
          onClick={() => onSelectShelf(shelf.id)}
          className="p-4 border rounded-lg shadow-sm hover:bg-gray-100"
        >
          {shelf.name}
        </button>
      ))}
    </div>
  );
}
