import { gql } from "urql";

export const ADD_EVENT = gql`
  mutation CreateEvent(
    $name: String!,
    $description: String,
    $startTime: Time!,
    $endTime: Time!
    $color: EventColor!) {
    createEvent(input: {
      name: $name,
      description: $description,
      startTime: $startTime,
      endTime: $endTime,
      color: $color
    }) {
      id
      name
      description
      startTime
      endTime
      color
    }
  }
`;

export const DELETE_EVENT = gql`
  mutation DeleteEvent($id: UUID!) {
    deleteEvent(id: $id)
  }
`;
