import { Client, fetchExchange, subscriptionExchange } from "urql";
import { createClient as createWSClient } from "graphql-ws";

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;
const WEBSOCKET_URL = import.meta.env.VITE_WEBSOCKET_URL;

const wsClient = createWSClient({ url: WEBSOCKET_URL });

export const client = new Client({
  url: BACKEND_URL,
  suspense: true,
  exchanges: [
    fetchExchange,
    subscriptionExchange({
      forwardSubscription(request) {
        const input = { ...request, query: request.query || "" };
        return {
          subscribe(sink) {
            const unsubscribe = wsClient.subscribe(input, sink);
            return { unsubscribe };
          }
        };
      }
    })
  ]
});
