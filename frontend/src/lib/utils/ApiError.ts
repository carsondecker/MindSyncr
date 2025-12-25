export class ApiError extends Error {
  readonly code: string;
  readonly status?: number;

  constructor(code: string, message: string) {
    super(message);
    this.code = code;

    Object.setPrototypeOf(this, new.target.prototype);
  }
}
