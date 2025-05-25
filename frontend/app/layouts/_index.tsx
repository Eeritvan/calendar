import { NavLink, Outlet } from "react-router";
import Sidebar from "~/components/sidebar";

const SidebarLayout = () => {
  return (
    <div className="flex flex-row h-dvh">
      <Sidebar>
        <NavLink to="/"> home </NavLink>
        <NavLink to="/test"> test </NavLink>
      </Sidebar>
      <Outlet />
    </div>
  );
};

export default SidebarLayout;
