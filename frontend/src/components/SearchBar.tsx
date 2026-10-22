import React from "react";
import { Search, X } from "lucide-react";

interface SearchBarProps {
  value: string;
  onChange: (val: string) => void;
  onClear: () => void;
}

export const SearchBar: React.FC<SearchBarProps> = ({ value, onChange, onClear }) => {
  return (
    <div style={{
      position: "relative",
      display: "flex",
      alignItems: "center",
      width: "100%",
      maxWidth: "520px"
    }}>
      <Search
        size={17}
        color="#94A3B8"
        strokeWidth={2}
        style={{ position: "absolute", left: "16px", pointerEvents: "none" }}
      />
      <input
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder="Search dosa, biryani, meals..."
        style={{
          width: "100%",
          padding: "11px 40px 11px 44px",
          backgroundColor: "#F4F5F7",
          border: "1px solid #EAECEF",
          borderRadius: "9999px",
          fontSize: "14px",
          color: "#171C24",
          outline: "none",
          transition: "all 0.2s ease"
        }}
        onFocus={(e) => {
          e.currentTarget.style.backgroundColor = "#ffffff";
          e.currentTarget.style.borderColor = "#FF5A1F";
          e.currentTarget.style.boxShadow = "0 0 0 3px rgba(255, 90, 31, 0.12)";
        }}
        onBlur={(e) => {
          e.currentTarget.style.backgroundColor = "#F4F5F7";
          e.currentTarget.style.borderColor = "#EAECEF";
          e.currentTarget.style.boxShadow = "none";
        }}
      />
      {value && (
        <button
          type="button"
          onClick={onClear}
          style={{
            position: "absolute",
            right: "14px",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            color: "#94A3B8",
            cursor: "pointer",
            padding: "2px"
          }}
          aria-label="Clear search"
        >
          <X size={16} />
        </button>
      )}
    </div>
  );
};
