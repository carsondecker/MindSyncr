export type ApiErrorCode = 
  | "BAD_REQUEST"
  | "USER_ALREADY_EXISTS"
  | "INVALID_CREDENTIALS"
  | "INVALID_REFRESH_TOKEN"
  | "INVALID_ACCESS_TOKEN"
  | "MISSING_ACCESS_TOKEN"
  | "FORBIDDEN"
  | "USER_NOT_FOUND"
  | "VALIDATION_FAIL"
  | "BAD_RESPONSE"
  | "DBTX_FAIL"
  | "MARSHAL_FAIL"
  | "HASH_FAIL"
  | "JWT_FAIL"
  | "REFRESH_FAIL"
  | "TX_BEGIN_FAIL"
  | "TX_COMMIT_FAIL"
  | "REFRESH_REVOKE_FAIL"
  | "GET_USER_DATA_FAIL"
  | "CREATE_JOIN_CODE_FAIL"
  | "BROADCAST_FAIL"
  | "NETWORK"
// TODO: remove broadcast fail when broadcasts become async

export class ApiError extends Error {
  readonly code: ApiErrorCode;
  readonly status: number;

  constructor(status: number, code: ApiErrorCode, message: string) {
    super(message);
    this.status = status;
    this.code = code;

    Object.setPrototypeOf(this, new.target.prototype);
  }
}
