import { Await, useParams, redirect } from "react-router";
import SingleDate from "./SingleDate";
import dayjs from "dayjs";
import type { Route } from "./+types";
import { client } from "~/api/graphql";
import { GET_EVENTS_BY_TIME_RANGE } from "../api/queries";
import { Suspense } from "react";
import type { Event } from "~/types";
import isBetween from "dayjs/plugin/isBetween";
import { urlDateSchema } from "../validation/dateUrl";

dayjs.extend(isBetween);

export const loader = ({ params }: Route.LoaderArgs) => {
  if (!params.startDate) {
    const currentDate = dayjs().format("YYYY-MM-DD");
    return redirect(`/week/${currentDate}`);
  }

  const dateParseResult = urlDateSchema.safeParse(params.startDate);
  if (!dateParseResult.success) {
    const currentDate = dayjs().format("YYYY-MM-DD");
    return redirect(`/week/${currentDate}`);
  }

  const startDate = params.startDate ? dayjs(params.startDate) : dayjs();

  const result = client.query(GET_EVENTS_BY_TIME_RANGE, {
    startTime: startDate,
    endTime: startDate.add(7, "day")
  }).toPromise();

  return { events: result };
};

const HOURS = Array.from(
  { length: 24 }, (_, i) => `${i.toString().padStart(2, "0")}:00`
);

const Week = ({ loaderData }: Route.ComponentProps) => {
  const { startDate } = useParams();
  const startDateObj = startDate ? dayjs(startDate) : dayjs();

  return (
    <div className={`grid grid-cols-[minmax(0,1fr)_repeat(7,_minmax(0,4fr))]
      grid-rows-[repeat(25,_40px)] grid-flow-col`}
    >
      <div />
      {HOURS.map((hour, index) => (
        <div key={index}>
          {hour}
        </div>
      ))}

      {Array.from({ length: 7 }, (_, index) => {
        const currentDate = startDateObj.add(index, "day");
        return (
          <div key={index} className="contents">
            <div>
              { currentDate.format("YYYY-MM-DD") }
            </div>
            <Suspense
              key={index}
              fallback={<SingleDate events={[]} />}
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

                  return <SingleDate events={dateEvents} />;
                }}
              </Await>
            </Suspense>
          </div>
        );
      })}
    </div>
  );
};

export default Week;
