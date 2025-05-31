import { gql } from "urql";

export const GET_QUERY = gql`
  query getAllEvents {
    allEvents {
      id
      name
      description
      startTime
      endTime
    }
  }
`;
