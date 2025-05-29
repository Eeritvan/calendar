import { useRef } from "react";
import { createCookie, data, NavLink, Outlet, useFetcher } from "react-router";
import type { Route } from "../../+types/root";

export const prefs = createCookie("prefs");

export async function loader({ request }: Route.LoaderArgs) {
  const cookieHeader = request.headers.get("Cookie");
  const cookie = (await prefs.parse(cookieHeader)) || {};
  return data({ sidebarWidth: cookie.sidebarWidth || 250 });
}

const Sidebar = ({ loaderData }: Route.ComponentProps) => {
  const fetcher = useFetcher();
  const sidebarRef = useRef<HTMLDivElement>(null);

  const { sidebarWidth = 250 } = loaderData  || {};

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
      const clampedWidth = Math.max(150, Math.min(500, newWidth));
      fetcher.submit(
        { sidebarWidth: clampedWidth },
        {
          method: "post",
          action: "/changeWidth"
        }
      );

      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleMouseUp);

      document.body.style.removeProperty("user-select");
      document.body.style.removeProperty("cursor");
    };

    document.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("mouseup", handleMouseUp);
  };

  return (
    <div className="flex flex-row h-dvh">
      <div
        ref={sidebarRef}
        className="bg-cyan-700"
        style={{ width: `${sidebarWidth}px` }}
      >
        <NavLink to="/"> home </NavLink>
        <NavLink to="/test"> test </NavLink>
      </div>
      <div
        className="w-2 bg-red-500"
        onMouseDown={startResizing}
      />
      <Outlet />
    </div>
  );
};

export default Sidebar;
