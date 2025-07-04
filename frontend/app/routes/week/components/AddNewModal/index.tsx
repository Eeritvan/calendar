import { useEffect, useRef } from "react";
import { useSearchParams, useNavigate } from "react-router";
import AddNewForm from "./AddNewForm";
import type { ColorHex } from "@/types";

const AddNewModal = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const modalRef = useRef<HTMLDivElement>(null);

  const startTime = searchParams.get("startTime") || "";
  const endTime = searchParams.get("endTime") || "";

  const closeModal = () => void navigate("..");

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        closeModal();
      }
    };
    document.addEventListener("keydown", handleKeyDown);

    return () => {
      document.removeEventListener("keydown", handleKeyDown);
    };
  });

  const handleColorChange = (color: ColorHex) => {
    if (modalRef.current) {
      modalRef.current.style.backgroundColor = color;
    }
  };

  return (
    <div
      ref={modalRef}
      className="fixed bg-gray-600"
    >
      <button onClick={closeModal}>
        x
      </button>
      <AddNewForm
        selectedTimeRange={{ startTime, endTime }}
        onColorChange={handleColorChange}
      />
    </div>
  );
};

export default AddNewModal;
