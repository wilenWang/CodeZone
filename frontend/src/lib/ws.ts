export type ServerEvent = {
  type: "message.created" | "conversation.updated" | "message.failed";
  payload: unknown;
};

export function wsUrl(origin: string, path: string): string {
  const url = new URL(path, origin);
  url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
  return url.toString();
}

export function connectEvents(
  token: string,
  onEvent: (event: ServerEvent) => void,
  onStatus: (connected: boolean) => void,
) {
  let closed = false;
  let socket: WebSocket | null = null;
  let reconnectTimer: number | null = null;

  function connect() {
    if (closed) {
      return;
    }
    socket = new WebSocket(`${wsUrl(window.location.origin, "/api/ws")}?access_token=${encodeURIComponent(token)}`);
    socket.onopen = () => onStatus(true);
    socket.onclose = () => {
      onStatus(false);
      if (!closed) {
        reconnectTimer = window.setTimeout(connect, 1200);
      }
    };
    socket.onmessage = (message) => {
      onEvent(JSON.parse(message.data) as ServerEvent);
    };
  }

  connect();

  return () => {
    closed = true;
    if (reconnectTimer !== null) {
      window.clearTimeout(reconnectTimer);
    }
    socket?.close();
  };
}
