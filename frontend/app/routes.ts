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
    route("/week/:startDate?", "routes/week/components/index.tsx"),
    route("/month/:date?", "routes/month/components/index.tsx"),
    route("test", "routes/test.tsx"),
    route("settings", "routes/settings/components/index.tsx")
  ]),
  route("changeSidebar", "routes/sidebar/actions/changeSidebar.tsx")
] satisfies RouteConfig;
