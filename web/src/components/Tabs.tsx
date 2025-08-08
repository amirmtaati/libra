"use client";
import { useState } from "react";

export default function Tabs({ tabs }) {
  const [activeTab, setActiveTab] = useState(Object.keys(tabs)[0]);

  return (
    <div>
      {/* Tab buttons */}
      <div className="flex border-b border-gray-300">
        {Object.keys(tabs).map((tab) => (
          <button
            key={tab}
            onClick={() => setActiveTab(tab)}
            className={`px-4 py-2 -mb-px border-b-2 ${
              activeTab === tab
                ? "border-blue-500 text-blue-500"
                : "border-transparent text-gray-500 hover:text-gray-700"
            }`}
          >
            {tab}
          </button>
        ))}
      </div>

      {/* Tab content */}
      <div className="mt-4">{tabs[activeTab]}</div>
    </div>
  );
}
