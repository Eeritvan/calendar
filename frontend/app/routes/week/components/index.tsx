import {
  Await,
  useParams,
  redirect,
  useNavigate,
  Outlet,
  useLocation
} from "react-router";
import SingleDate from "./SingleDate";
import dayjs from "dayjs";
import type { Route } from "./+types";
import { client } from "@/api/graphql";
import { GET_EVENTS_BY_TIME_RANGE } from "../api/queries";
import { Suspense, useEffect, useState } from "react";
import type { Event } from "@/types";
import isBetween from "dayjs/plugin/isBetween";
import { urlDateSchema } from "../validation/dateUrl";
import HourColumn from "./HourColumn";

// eslint-disable-next-line import-x/no-named-as-default-member
dayjs.extend(isBetween);

interface GetEventsResponse {
  data?: {
    eventsByTimeRange?: Event[];
  };
}

export const meta = () => {
  return [
    { title: "week view" },
    { name: "weeks", content: "Testing react router!" }
  ];
};

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

const Week = ({ loaderData }: Route.ComponentProps) => {
  const { startDate } = useParams();
  const navigate = useNavigate();
  const location = useLocation();

  const startDateObj = startDate ? dayjs(startDate) : dayjs();
  const [selectedTimeRange, setSelectedTimeRange] = useState<{
    startTime: string;
    endTime: string;
  }>({ startTime: "", endTime: "" });
  const emptyEvents: Event[] = [];

  useEffect(() => {
    if (!location.pathname.endsWith("/add")) {
      setSelectedTimeRange({ startTime: "", endTime: "" });
    }
  }, [location]);

  const handleSelect = (startTime: string, endTime: string) => {
    void navigate(`add?startTime=${startTime}&endTime=${endTime}`);
    setSelectedTimeRange({ startTime, endTime });
  };

  return (
    <div className={`grid grid-cols-[minmax(0,1fr)_repeat(7,_minmax(0,4fr))]
      grid-rows-[repeat(25,_50px)] grid-flow-col`}
    >
      <HourColumn />

      {Array.from({ length: 7 }, (_, index) => {
        const currentDate = startDateObj.add(index, "day");
        return (
          <div key={index} className="contents">
            <div className="border-2">
              { currentDate.format("YYYY-MM-DD") }
            </div>
            <Suspense
              key={index}
              fallback={
                <SingleDate
                  date={currentDate}
                  events={emptyEvents}
                  handleSelect={handleSelect}
                />
              }
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
                    <SingleDate
                      date={currentDate}
                      events={dateEvents}
                      handleSelect={handleSelect}
                      showSelectedTime={
                        selectedTimeRange.startTime &&
                        dayjs(currentDate)
                          .isSame(dayjs(selectedTimeRange.startTime), "day")}
                      selectedTimeRange={selectedTimeRange}
                    />
                  );
                }}
              </Await>
            </Suspense>
          </div>
        );
      })}
      <Outlet />
    </div>
  );
};

export default Week;
