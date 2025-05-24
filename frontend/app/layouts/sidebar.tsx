import { NavLink, Outlet } from "react-router";

const SidebarLayout = () => {
  return (
    <>
      <NavLink to="/"> home </NavLink>
      <Outlet />
    </>
  );
};

export default SidebarLayout;
