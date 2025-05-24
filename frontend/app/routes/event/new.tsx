import { gql } from "urql";
import { client } from "~/root";
import type { Route } from "./+types/new";
import { redirect } from "react-router";

const ADD_QUERY = gql`
  mutation CreateEvent(
    $name: String!,
    $description: String,
    $startTime: Time!,
    $endTime: Time!) {
    createEvent(input: {
      name: $name,
      description: $description,
      startTime: $startTime,
      endTime: $endTime,
    }) {
      id
      name
      description
      startTime
      endTime
    }
  }
`;

export const clientAction = async ({ request }: Route.ActionArgs) => {
  const formData = await request.formData();

  const name = formData.get("name") as string;
  const description = formData.get("description") as string;
  const time1 = new Date();
  const time2 = new Date();

  await client.mutation(
    ADD_QUERY, {
      name: name,
      description: description,
      startTime: time1,
      endTime: time2
    }).toPromise();

  return redirect("/");
};
