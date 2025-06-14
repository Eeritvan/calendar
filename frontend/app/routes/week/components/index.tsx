import { Await, useParams } from "react-router";
import SingleDate from "./SingleDate";
import dayjs from "dayjs";
import type { Route } from "./+types";
import { client } from "~/api/graphql";
import { GET_EVENTS_BY_TIME_RANGE } from "../api/queries";
import { Suspense } from "react";
import type { Event } from "~/types";
import isBetween from "dayjs/plugin/isBetween";

dayjs.extend(isBetween);

export const loader = ({ params }: Route.LoaderArgs) => {
  const startDate = params.startDate ? dayjs(params.startDate) : dayjs();

  const result = client.query(GET_EVENTS_BY_TIME_RANGE, {
    startTime: startDate,
    endTime: startDate.add(7, "day")
  }).toPromise();

  return { events: result };
};

const Week = ({ loaderData }: Route.ComponentProps) => {
  const { startDate } = useParams();
  const startDateObj = startDate ? dayjs(startDate) : dayjs();

  return (
    <div className="flex w-full">
      {Array.from({ length: 7 }, (_, index) => {
        const currentDate = startDateObj.add(index, "day");

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

export default Week;
