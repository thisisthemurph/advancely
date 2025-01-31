class ApiError extends Error {
  constructor(public status: number, public statusText: string, message?: string) {
    super(message);
    this.name = "ApiError";
  }
}

export default ApiError;
