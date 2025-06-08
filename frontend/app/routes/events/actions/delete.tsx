import { client } from "~/api/graphql";
import { redirect } from "react-router";
import { DELETE_EVENT } from "../api/mutations";
import type { Route } from "./+types/delete";

export const action = async ({ params }: Route.ActionArgs) => {
  const id = params.eventId;
  await client.mutation(DELETE_EVENT, { id: id }).toPromise();
  return redirect("/");
};
