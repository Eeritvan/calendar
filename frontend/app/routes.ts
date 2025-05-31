import {
  type RouteConfig,
  layout,
  prefix,
  route
} from "@react-router/dev/routes";

export default [
  layout("routes/sidebar/components/index.tsx", [
    route("/", "routes/events/components/index.tsx", [
      ...prefix("event", [
        route("new", "routes/events/actions/new.tsx"),
        route("delete/:eventId", "routes/events/actions/delete.tsx")
      ])
    ]),
    route("test", "routes/test.tsx")
  ]),
  route("changeSidebar", "routes/sidebar/actions/changeSidebar.tsx")
] satisfies RouteConfig;
