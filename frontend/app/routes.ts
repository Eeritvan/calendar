import {
  type RouteConfig,
  index,
  layout,
  route
} from "@react-router/dev/routes";

export default [
  layout("layouts/sidebar.tsx", [
    index("routes/mainView/index.tsx")
  ]),
  route("new", "routes/event/new.tsx"),
  route("delete/:EventID", "routes/event/delete.tsx")
] satisfies RouteConfig;
