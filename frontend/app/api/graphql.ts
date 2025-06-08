import { Client, fetchExchange, subscriptionExchange } from "urql";
import { createClient as createWSClient } from "graphql-ws";

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

const wsClient = createWSClient({
  url: BACKEND_URL.replace("http://", "ws://") || "ws://localhost:8081/api"
});

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
