import React, { useState, useRef, useEffect } from "react";
import { ChevronDown, Check } from "lucide-react";

export const NEIGHBOURHOODS = [
  "Indiranagar",
  "Koramangala",
  "HSR Layout",
  "BTM Layout",
  "Whitefield",
  "Marathahalli",
  "Malleshwaram",
  "Jayanagar"
];

interface LocationSelectorProps {
  selected: string;
  onSelect: (neighbourhood: string) => void;
}

export const LocationSelector: React.FC<LocationSelectorProps> = ({ selected, onSelect }) => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const displayLocation = selected === "all" ? "Indiranagar" : selected;

  return (
    <div style={{ position: "relative" }} ref={dropdownRef}>
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        style={{
          display: "flex",
          alignItems: "center",
          gap: "10px",
          padding: "7px 16px 7px 12px",
          borderRadius: "9999px",
          backgroundColor: "#F4F5F7",
          border: "1px solid #EAECEF",
          color: "#171C24",
          fontSize: "13px",
          cursor: "pointer",
          transition: "all 0.2s ease"
        }}
        onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = "#ECEEEF")}
        onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = "#F4F5F7")}
      >
        {/* Solid Orange Map Pin with White Center Dot */}
        <div style={{ display: "flex", alignItems: "center", justifyContent: "center", flexShrink: 0 }}>
          <svg width="18" height="20" viewBox="0 0 24 28" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 0C5.372 0 0 5.372 0 12c0 9 12 16 12 16s12-7 12-16c0-6.628-5.372-12-12-12z" fill="#FF5A1F"/>
            <circle cx="12" cy="11" r="4.5" fill="#ffffff"/>
          </svg>
        </div>

        <div style={{ textAlign: "left", lineHeight: 1.25 }}>
          <div style={{ fontSize: "11px", color: "#64748B", fontWeight: 500 }}>
            Delivering to
          </div>
          <div style={{ fontWeight: 700, fontSize: "13.5px", color: "#171C24" }}>
            {displayLocation}, Bengaluru
          </div>
        </div>

        <ChevronDown size={15} color="#475569" style={{ marginLeft: "4px" }} />
      </button>

      {isOpen && (
        <div
          style={{
            position: "absolute",
            top: "calc(100% + 8px)",
            left: 0,
            width: "250px",
            backgroundColor: "#ffffff",
            borderRadius: "14px",
            boxShadow: "0 10px 30px rgba(0,0,0,0.12)",
            border: "1px solid #E2E8F0",
            zIndex: 100,
            padding: "8px 0"
          }}
        >
          <div style={{ padding: "8px 16px", fontSize: "11px", fontWeight: 800, color: "#94A3B8", textTransform: "uppercase", letterSpacing: "0.5px" }}>
            Select Bengaluru Location
          </div>
          <button
            type="button"
            onClick={() => { onSelect("all"); setIsOpen(false); }}
            style={{
              width: "100%",
              display: "flex",
              alignItems: "center",
              justifyContent: "space-between",
              padding: "10px 16px",
              fontSize: "13.5px",
              textAlign: "left",
              backgroundColor: selected === "all" ? "#FFF0E6" : "transparent",
              color: selected === "all" ? "#FF5A1F" : "#171C24",
              fontWeight: selected === "all" ? 700 : 500,
              cursor: "pointer"
            }}
          >
            <span>All Bengaluru Areas</span>
            {selected === "all" && <Check size={15} />}
          </button>
          <div style={{ height: "1px", backgroundColor: "#F1F5F9", margin: "4px 0" }} />
          {NEIGHBOURHOODS.map(area => (
            <button
              key={area}
              type="button"
              onClick={() => { onSelect(area); setIsOpen(false); }}
              style={{
                width: "100%",
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between",
                padding: "9px 16px",
                fontSize: "13.5px",
                textAlign: "left",
                backgroundColor: selected === area ? "#FFF0E6" : "transparent",
                color: selected === area ? "#FF5A1F" : "#171C24",
                fontWeight: selected === area ? 700 : 500,
                cursor: "pointer"
              }}
            >
              <span>{area}</span>
              {selected === area && <Check size={15} />}
            </button>
          ))}
        </div>
      )}
    </div>
  );
};
