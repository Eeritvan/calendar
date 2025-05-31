import { gql } from "urql";

export const ADD_EVENT = gql`
  mutation CreateEvent(
    $name: String!,
    $description: String,
    $startTime: Time!,
    $endTime: Time!) {
    createEvent(input: {
      name: $name,
      description: $description,
      startTime: $startTime,
      endTime: $endTime,
    }) {
      id
      name
      description
      startTime
      endTime
    }
  }
`;

export const DELETE_EVENT = gql`
  mutation DeleteEvent($id: UUID!) {
    deleteEvent(id: $id)
  }
`;
