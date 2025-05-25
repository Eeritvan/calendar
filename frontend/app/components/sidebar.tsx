import { type ReactNode, useState, useRef } from "react";

const Sidebar = ({ children }: { children: ReactNode }) => {
  const [size, setSize] = useState<number>(400);
  const sidebarRef = useRef<HTMLDivElement>(null);

  const startResizing = () => {
    document.body.style.userSelect = "none";
    document.body.style.cursor = "col-resize";

    const handleMouseMove = (moveEvent: MouseEvent) => {
      const newWidth = moveEvent.clientX;
      if (newWidth >= 150 && newWidth <= 500 && sidebarRef.current) {
        sidebarRef.current.style.width = `${newWidth}px`;
      }
    };

    const handleMouseUp = (moveEvent: MouseEvent) => {
      const newWidth = moveEvent.clientX;
      if (newWidth >= 150 && newWidth <= 500) {
        setSize(newWidth);
      }

      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleMouseUp);

      document.body.style.removeProperty("user-select");
      document.body.style.removeProperty("cursor");
    };

    document.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("mouseup", handleMouseUp);
  };

  return (
    <>
      <div
        ref={sidebarRef}
        className="bg-cyan-700"
        style={{ width: `${size}px` }}
      >
        {children}
      </div>
      <div
        className="w-2 bg-red-500"
        onMouseDown={startResizing}
      />
    </>
  );
};

export default Sidebar;
