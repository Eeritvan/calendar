import { gql } from "urql";

export const SUBS_TEST = gql`
  subscription EventChanged {
    eventChanged {
      action
      event {
        id
        name
        description
        startTime
        endTime
      }
    }
  }
`;
