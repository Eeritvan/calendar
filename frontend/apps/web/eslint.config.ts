import globals from "globals";
import tseslint from "typescript-eslint";
import { baseConfig } from "@repo/eslint-config/base";
import { defineConfig } from "eslint/config";

export default defineConfig([
   baseConfig,
   ...tseslint.configs.recommended,
  {
    languageOptions: { globals: globals.browser },
    ignores: [
      "app.config.ts"
    ]
  }
]);
