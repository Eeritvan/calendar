import dayjs from "dayjs";
import { Await, redirect, useParams } from "react-router";
import type { Route } from "./+types";
import { Suspense } from "react";
import SingleDate from "./SingleDate";
import { urlDateSchema } from "../validation/dateUrl";
import { client } from "@/api/graphql";
import { GET_EVENTS_BY_TIME_RANGE } from "../api/query";
import type { Event } from "@/types";
import isBetween from "dayjs/plugin/isBetween";

// eslint-disable-next-line import-x/no-named-as-default-member
dayjs.extend(isBetween);

interface GetEventsResponse {
  data?: {
    eventsByTimeRange?: Event[];
  };
}

export const meta = () => {
  return [
    { title: "month view" },
    { name: "months", content: "Testing react router!" }
  ];
};

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
  const daysInMonth = dayjs(startOfMonth).daysInMonth();

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

  const emptyEvents: Event[] = [];

  return (
    <div className="grid grid-cols-7 m-2 h-dvh">
      {Array.from({ length: daysInMonth }, (_, index) => {
        const currentDate = parsedDate.add(index, "day");

        return (
          <Suspense
            key={currentDate.format("YYYY-MM-DD")}
            fallback={ <SingleDate date={currentDate} events={emptyEvents} /> }
          >
            <Await resolve={loaderData.events}>
              {(data: GetEventsResponse) => {
                const events: Event[] = data.data?.eventsByTimeRange || [];
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
