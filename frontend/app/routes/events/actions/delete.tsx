import { client } from "~/api/graphql";
import { redirect } from "react-router";
import { DELETE_QUERY } from "../api/mutations";
import type { Route } from "./+types/delete";

export const clientAction = async ({ params }: Route.ActionArgs) => {
  const id = params.eventId;
  await client.mutation(DELETE_QUERY, { id: id }).toPromise();
  return redirect("/");
};
