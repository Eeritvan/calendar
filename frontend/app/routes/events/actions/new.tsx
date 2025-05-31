import { client } from "~/api/graphql";
import { redirect } from "react-router";
import dayjs from "dayjs";
import { eventValidationSchema } from "../validation/schemas";
import { ADD_QUERY } from "../api/mutations";
import type { Time } from "../types";
import type { Route } from "./+types/new";

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
