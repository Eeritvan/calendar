import { z } from "zod/v4-mini";

const yyyyMmFormat = /^\d{4}-(0[1-9]|1[0-2])$/;
export const urlDateSchema = z.string().check(z.regex(yyyyMmFormat));
