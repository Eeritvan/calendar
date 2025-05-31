import { useRef, useState } from "react";
import { createCookie, data, NavLink, Outlet, useFetcher } from "react-router";
import type { Route } from "./+types";

export const prefs = createCookie("prefs");

export async function loader({ request }: Route.LoaderArgs) {
  const cookieHeader = request.headers.get("Cookie");
  const cookie = (await prefs.parse(cookieHeader)) || {};
  return data({
    sidebarWidth: cookie.sidebarWidth || 250,
    isCollapsed: cookie.isCollapsed || false
  });
}

interface SidebarProps {
  loaderData: {
    isCollapsed?: boolean;
    sidebarWidth?: number;
  };
}

const Sidebar = ({ loaderData }: SidebarProps) => {
  const fetcher = useFetcher();
  const sidebarRef = useRef<HTMLDivElement>(null);
  const [isCollapsed, setIsCollapsed] =
    useState<boolean>(loaderData?.isCollapsed || false);

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
          method: "POST",
          action: "/changeSidebar"
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

  const toggleSidebar = () => {
    setIsCollapsed(!isCollapsed);
    fetcher.submit(
      { isCollapsed: !isCollapsed },
      {
        method: "POST",
        action: "/changeSidebar"
      }
    );
  };

  return (
    <>
      { isCollapsed ?
        <>
          <button onClick={toggleSidebar}>
            toggle
          </button>
          <Outlet />
        </> :
        <div className="flex flex-row h-dvh">
          <div
            ref={sidebarRef}
            className="bg-cyan-700 relative"
            style={{ width: `${sidebarWidth}px` }}
          >
            <NavLink to="/"> home </NavLink>
            <NavLink to="/test"> test </NavLink>
            <NavLink to="#settings"> settings </NavLink>
            <button
              onClick={toggleSidebar}
              className="absolute top-0 right-0"
            >
              toggle
            </button>
          </div>
          <div
            className="w-2 bg-red-500"
            onMouseDown={startResizing}
          />
          <Outlet />
        </div>
      }
    </>
  );
};

export default Sidebar;
