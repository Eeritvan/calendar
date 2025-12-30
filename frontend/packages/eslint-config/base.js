import js from "@eslint/js";
import turboPlugin from "eslint-plugin-turbo";
import onlyWarn from "eslint-plugin-only-warn";
import stylistic from '@stylistic/eslint-plugin'

export const baseConfig = [
  js.configs.recommended,
  {
    files: ["**/*.{js,mjs,cjs,ts,mts,cts,jsx,tsx}"],
    plugins: {
      turbo: turboPlugin,
      '@stylistic': stylistic,
    },
    rules: {
      "turbo/no-undeclared-env-vars": "warn",
      semi: ["error", "always"],
      quotes: ["error", "double"],
      eqeqeq: ["error"],
      camelcase: ["error", { "properties": "always" }],
      '@stylistic/indent': ['error', 2],
      "no-trailing-spaces": ["error"],
      "linebreak-style": ["error", "unix"],
      "arrow-spacing": ["error", { "before": true, "after": true }],
      "object-curly-spacing": ["error", "always"],
      "max-len": ["error", 80],
      "no-console": ["warn"],
      "no-multiple-empty-lines": ["error", { "max": 1 }],
      "eol-last": ["error", "always"],
      "comma-dangle": ["error", "always"],
      "max-depth": ["error", 2],
      "no-else-return": ["error"],
      "comma-spacing": ["error", { "before": false, "after": true }],
      "no-var": ["error"],
      "prefer-const": ["error"],
    },
    ignores: [
      "dist/**",
      "eslint.config.*",
      "node_modules/**"
    ],
  },
  {
    plugins: {
      onlyWarn,
    },
  },
];
