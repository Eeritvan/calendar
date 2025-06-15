import { gql } from "urql";

// todo: better name for this idk
export const GET_EVENTS_BY_TIME_RANGE = gql`
  query eventsByTimeRange(
    $startTime: Time!
    $endTime: Time!
  ) {
    eventsByTimeRange(
      startTime: $startTime
      endTime: $endTime
    ) {
      id
      name
      description
      startTime
      endTime
    }
  }
`;
