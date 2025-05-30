import { gql } from "urql";
import { client } from "~/root";
import type { Route } from "./+types/new";
import { redirect } from "react-router";
import dayjs from "dayjs";
import { z } from "zod/v4-mini";

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

const eventValidationSchema = z.object({
  name: z.string().check(z.minLength(3), z.maxLength(100), z.trim()),
  description: z.optional(z.string().check(z.maxLength(1000), z.trim())),
  startTime: z.string().check(z.minLength(1, "Start time is required")),
  endTime: z.string().check(z.minLength(1, "missing stuff"))
}).check(z.refine((data) => {
  const start = new Date(data.startTime);
  const end = new Date(data.endTime);
  return end > start;
}, {
  message: "End time must be after start time",
  path: ["endTime"]
}));

export const clientAction = async ({ request }: Route.ActionArgs) => {
  const formData = await request.formData();

  const rawData = {
    name: formData.get("name") as string,
    description: formData.get("description") as string,
    startTime: dayjs(formData.get("startTime") as string).format() as Time,
    endTime: dayjs(formData.get("endTime") as string).format() as Time
  };

  const result = eventValidationSchema.safeParse(rawData);

  if (!result.success) {
    return {
      errors: result.error.issues,
      data: rawData
    };
  }

  await client.mutation(ADD_QUERY, rawData).toPromise();

  return redirect("/");
};
