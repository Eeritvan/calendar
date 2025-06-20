import dayjs from "dayjs";
import { Await, redirect, useParams } from "react-router";
import type { Route } from "./+types";
import { Suspense } from "react";
import SingleDate from "~/routes/month/components/SingleDate";
import { urlDateSchema } from "../validation/date";
import { client } from "~/api/graphql";
import { GET_EVENTS_BY_TIME_RANGE } from "../api/query";
import type { Event } from "~/types";
import isBetween from "dayjs/plugin/isBetween";

dayjs.extend(isBetween);

export const loader = ({ params }: Route.LoaderArgs) => {
  if (!params.date) {
    const currentDate = dayjs().format("YYYY-MM");
    return redirect(`/month/${currentDate}`);
  }

  const dateParseResult = urlDateSchema.safeParse(params.date);
  if (!dateParseResult.success) {
    const currentDate = dayjs().format("YYYY-MM");
    return redirect(`/month/${currentDate}`);
  }

  const startOfMonth = dayjs(params.date);
  const daysInMonth = startOfMonth ? dayjs(startOfMonth).daysInMonth() : 0;

  const result = client.query(GET_EVENTS_BY_TIME_RANGE, {
    startTime: startOfMonth,
    endTime: startOfMonth.add(daysInMonth, "day")
  }).toPromise();

  return { events: result };
};

const Month = ({ loaderData }: Route.ComponentProps) => {
  const { date } = useParams();
  const parsedDate = dayjs(date);
  const daysInMonth = parsedDate.isValid() ? parsedDate.daysInMonth() : 0;

  return (
    <div className="grid grid-cols-7 m-2 h-dvh">
      {Array.from({ length: daysInMonth }, (_, index) => {
        const currentDate = parsedDate.add(index, "day");

        return (
          <Suspense
            key={index}
            fallback={ <SingleDate date={currentDate} events={[]} /> }
          >
            <Await resolve={loaderData.events}>
              {(data) => {
                const events: Event[] = data?.data?.eventsByTimeRange || [];
                const dateEvents: Event[] = events.filter((event: Event) => {
                  return dayjs(currentDate).isBetween(
                    dayjs(event.startTime),
                    dayjs(event.endTime),
                    "day",
                    "[]"
                  );
                });

                return (
                  <SingleDate date={currentDate} events={dateEvents} />
                );
              }}
            </Await>
          </Suspense>
        );
      })}
    </div>
  );
};

export default Month;
