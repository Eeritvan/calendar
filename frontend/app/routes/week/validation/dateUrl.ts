import { z } from "zod/v4-mini";

const yyyyMmDdFormat = /^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$/;
export const urlDateSchema = z.string().check(z.regex(yyyyMmDdFormat));
