import { data } from "react-router";
import type { Route } from "./+types/sidebar";
import { prefs } from "./sidebar";

export async function action({ request }: Route.ActionArgs) {
  const cookieHeader = request.headers.get("Cookie");
  const cookie = (await prefs.parse(cookieHeader)) || {};
  const formData = await request.formData();

  const sidebarWidth = Number(formData.get("sidebarWidth"));

  cookie.sidebarWidth = sidebarWidth;

  return data(sidebarWidth, {
    headers: {
      "Set-Cookie": await prefs.serialize(cookie)
    }
  });
}
