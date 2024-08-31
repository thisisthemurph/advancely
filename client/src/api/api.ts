import { z } from "zod";

export const ErrorResponseSchema = z.object({
  message: z.string(),
});

export type ErrorResponse = z.infer<typeof ErrorResponseSchema>;

export const handleErrorResponse = async (
  resp: Response,
  defaultError?: string
): Promise<ErrorResponse> => {
  if (!defaultError) {
    defaultError = "Unknown error processing request.";
  }

  const body = await resp.json();
  const error = ErrorResponseSchema.safeParse(body);

  return error.success && !isDefaultError(resp.status, error.data)
    ? error.data
    : { message: defaultError };
};

function isDefaultError(status: number, error: ErrorResponse): boolean {
  console.log({ status, error: error.message });

  // Define default error messages based on status codes
  const defaultErrors: { [key: number]: string } = {
    400: "Bad Request",
    401: "Unauthorized",
    403: "Forbidden",
    404: "Not Found",
    500: "Internal Server Error",
    502: "Bad Gateway",
    503: "Service Unavailable",
    504: "Gateway Timeout",
  };

  // Check if the status code exists in the defaultErrors object
  if (defaultErrors[status] && error.message === defaultErrors[status]) {
    return true;
  }

  return false;
}
