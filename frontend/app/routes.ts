import {
  type RouteConfig,
  index,
  layout,
  route
} from "@react-router/dev/routes";

export default [
  layout("routes/sidebar/sidebar.tsx", [
    index("routes/mainView/index.tsx"),
    route("test", "routes/test.tsx")
  ]),
  route("new", "routes/event/new.tsx"),
  route("delete/:EventID", "routes/event/delete.tsx"),
  route("changeSidebar", "routes/sidebar/changeSidebar.tsx")
] satisfies RouteConfig;
