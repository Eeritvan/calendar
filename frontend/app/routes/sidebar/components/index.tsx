import { useRef, useState } from "react";
import { createCookie, data, NavLink, Outlet, useFetcher } from "react-router";
import type { Route } from "./+types";

export const prefs = createCookie("prefs");

interface CookieProps {
  sidebarWidth?: number;
  isCollapsed?: boolean;
}

export async function loader({ request }: Route.LoaderArgs) {
  const cookieHeader = request.headers.get("Cookie");
  const cookie = await prefs.parse(cookieHeader) as CookieProps || {};
  return data({
    sidebarWidth: cookie.sidebarWidth || 250,
    isCollapsed: cookie.isCollapsed || false
  });
}

interface SidebarProps {
  loaderData: {
    isCollapsed: boolean;
    sidebarWidth: number;
  };
}

const Sidebar = ({ loaderData }: SidebarProps) => {
  const fetcher = useFetcher();
  const sidebarRef = useRef<HTMLDivElement>(null);
  const [isCollapsed, setIsCollapsed] =
    useState<boolean>(loaderData.isCollapsed || false);

  const { sidebarWidth = 250 } = loaderData;

  const startResizing = () => {
    document.body.style.userSelect = "none";
    document.body.style.cursor = "col-resize";

    const handleMouseMove = (moveEvent: MouseEvent) => {
      const newWidth = moveEvent.clientX;
      if (newWidth >= 150 && newWidth <= 500 && sidebarRef.current) {
        sidebarRef.current.style.width = `${newWidth.toString()}px`;
      }
    };

    const handleMouseUp = (moveEvent: MouseEvent) => {
      const newWidth = moveEvent.clientX;
      const clampedWidth = Math.max(150, Math.min(500, newWidth));
      void fetcher.submit(
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
    void fetcher.submit(
      { isCollapsed: !isCollapsed },
      {
        method: "POST",
        action: "/changeSidebar"
      }
    );
  };

  const widthStyle = { width: `${sidebarWidth.toString()}px` };

  return (
    <>
      { isCollapsed ?
        <>
          <button onClick={toggleSidebar}>
            toggle
          </button>
          <Outlet />
        </> :
        <div className="flex h-dvh">
          <div
            ref={sidebarRef}
            className="bg-cyan-700 relative"
            style={widthStyle}
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
            role="presentation"
            className="w-2 bg-red-500 flex-shrink-0"
            onMouseDown={startResizing}
          />
          <div className="flex-1 overflow-y-auto">
            <Outlet />
          </div>
        </div>
      }
    </>
  );
};

export default Sidebar;
