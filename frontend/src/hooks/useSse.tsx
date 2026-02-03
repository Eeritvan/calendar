import { useEffect } from "react"
import { API_URL } from "@/constants";

export const useSse = (connId: string) => {
  useEffect(() => {
    if (!connId) return;
    const eventSource = new EventSource(
      `${API_URL}/sse?stream=${connId}`,
      { withCredentials: true }
    );

    eventSource.addEventListener("open", () => {
      console.log("opened")
    })

    eventSource.onerror = (e) => {
      console.log(e);
    };

    // calendars
    eventSource.addEventListener("calendar/post", (e) => {
      const data = JSON.parse((e).data)
      console.log("calendar/post", data);
    });

    eventSource.addEventListener("calendar/edit", (e) => {
      const data = JSON.parse((e).data)
      console.log("calendar/edit", data);
    });

    eventSource.addEventListener("calendar/delete", (e) => {
      const data = (e).data
      console.log("calendar/delete", data);
    });

    // events
    eventSource.addEventListener("event/post", (e) => {
      const data = JSON.parse((e).data)
      console.log("event/post", data);
    });

    eventSource.addEventListener("event/edit", (e) => {
      const data = JSON.parse((e).data)
      console.log("event/edit", data);
    });

    eventSource.addEventListener("event/delete", (e) => {
      const data = (e).data
      console.log("event/delete", data);
    });

    return () => {
      eventSource.close();
    };
  }, [connId])
}
