const SHANGHAI_OFFSET_MS = 8 * 60 * 60 * 1000;

function toShanghaiDate(timestamp: number) {
  return new Date(timestamp * 1000 + SHANGHAI_OFFSET_MS);
}

function isValidTimestamp(timestamp: number | null | undefined): timestamp is number {
  return typeof timestamp === "number" && Number.isFinite(timestamp);
}

function pad(value: number) {
  return String(value).padStart(2, "0");
}

export function formatDate(timestamp: number | null | undefined) {
  if (!isValidTimestamp(timestamp)) return "";

  const date = toShanghaiDate(timestamp);
  return `${date.getUTCFullYear()}/${date.getUTCMonth() + 1}/${date.getUTCDate()}`;
}

export function formatDateTime(timestamp: number | null | undefined) {
  if (!isValidTimestamp(timestamp)) return "";

  const date = toShanghaiDate(timestamp);
  return `${formatDate(timestamp)} ${pad(date.getUTCHours())}:${pad(
    date.getUTCMinutes()
  )}:${pad(date.getUTCSeconds())}`;
}
