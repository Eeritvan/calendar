import js from "@eslint/js";
import globals from "globals";
import tseslint from "typescript-eslint";
import pluginReact from "eslint-plugin-react";
import { defineConfig } from "eslint/config";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import depend from "eslint-plugin-depend";

export default defineConfig([
  {
    files: ["**/*.{js,mjs,cjs,ts,jsx,tsx}"],
    plugins: {
      depend
    },
    extends: ["depend/flat/recommended"],
    languageOptions:{
      ecmaVersion: "latest",
      globals: globals.browser,
      parserOptions: {
        ecmaFeatures: {
          jsx: true
        }
      }
    }
  },

  js.configs.recommended,
  tseslint.configs.recommended,
  pluginReact.configs.flat.recommended,
  reactHooks.configs.recommended,
  reactRefresh.configs.recommended,

  { ignores: ["build", ".react-router", "eslint.config.js"] },
  {
    settings: {
      react: {
        version: "detect"
      }
    },
    rules: {
      "react-hooks/react-compiler": ["error"],
      "react-refresh/only-export-components": ["off"],
      semi: ["error", "always"],
      quotes: ["error", "double"],
      indent: ["error", 2],
      eqeqeq: ["error"],
      camelcase: ["error", { "properties": "always" }],
      "no-trailing-spaces": ["error"],
      "linebreak-style": ["error", "unix"],
      "arrow-spacing": ["error", { "before": true, "after": true }],
      "object-curly-spacing": ["error", "always"],
      "max-len": ["error", 80],
      "no-console": ["error"],
      "no-duplicate-imports": ["error"],
      "no-multiple-empty-lines": ["error", { "max": 1 }],
      "eol-last": ["error", "always"],
      "comma-dangle": ["error", "never"],
      "max-depth": ["error", 2],
      "no-else-return": ["error"],
      "comma-spacing": ["error", { "before": false, "after": true }],
      "no-var": ["error"],
      "prefer-const": ["error"],
      "react/jsx-closing-bracket-location": ["error"],
      "react/prefer-stateless-function": ["error"],
      "react/no-multi-comp": ["error"],
      "react/sort-comp": ["error"],
      "react/self-closing-comp": ["error"],
      "react/jsx-wrap-multilines": ["error"],
      "react/react-in-jsx-scope": ["off"]
    }
  }
]);
