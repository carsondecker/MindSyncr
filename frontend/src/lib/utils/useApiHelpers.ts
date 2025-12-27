import { ApiError } from "./ApiError";

export function canRetry(err: unknown) {
  if (!(err instanceof ApiError)) return false
  if (err.code === "NETWORK") return true
  if (err.status >= 500) return true
  return false
}