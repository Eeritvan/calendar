import { gql } from "urql";
import { client } from "~/root";
import type { Route } from "./+types/new";
import { redirect } from "react-router";
import dayjs from "dayjs";

type Time = string & { readonly __brand: unique symbol };

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
  const startTime = dayjs(formData.get("startTime") as string).format() as Time;
  const endTime = dayjs(formData.get("endTime") as string).format() as Time;

  await client.mutation(
    ADD_QUERY, {
      name: name,
      description: description,
      startTime: startTime,
      endTime: endTime
    }).toPromise();

  return redirect("/");
};
