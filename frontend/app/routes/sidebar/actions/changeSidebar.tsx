import { data } from "react-router";
import type { Route } from "./+types/changeSidebar";
import { prefs } from "../components";

interface CookieProps {
  sidebarWidth: number;
  isCollapsed: boolean;
}

export async function action({ request }: Route.ActionArgs) {
  const cookieHeader = request.headers.get("Cookie");
  const cookie = await prefs.parse(cookieHeader) as CookieProps || {};
  const formData = await request.formData();

  const sidebarWidth = formData.get("sidebarWidth");
  const isCollapsed = formData.get("isCollapsed");

  if (sidebarWidth !== null) {
    cookie.sidebarWidth = Number(sidebarWidth);
  }

  if (isCollapsed !== null) {
    cookie.isCollapsed = isCollapsed === "true";
  }

  return data(
    { sidebarWidth: cookie.sidebarWidth, isCollapsed: cookie.isCollapsed }, {
      headers: {
        "Set-Cookie": await prefs.serialize(cookie)
      }
    }
  );
}
