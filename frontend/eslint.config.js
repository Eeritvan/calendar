import js from "@eslint/js";
import globals from "globals";
import tseslint from "typescript-eslint";
import pluginReact from "eslint-plugin-react";
import { defineConfig } from "eslint/config";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import depend from "eslint-plugin-depend";
import jsxA11y from "eslint-plugin-jsx-a11y";
import barrelfiles from "eslint-plugin-barrel-files";
import reactCompiler from "eslint-plugin-react-compiler"
import { importX } from "eslint-plugin-import-x";
import tsParser from "@typescript-eslint/parser";

export default defineConfig([
  js.configs.recommended,
  tseslint.configs.strictTypeChecked,
  pluginReact.configs.flat.recommended,
  reactHooks.configs["recommended-latest"],
  reactRefresh.configs.vite,
  jsxA11y.flatConfigs.strict,
  barrelfiles.configs.recommended,
  reactCompiler.configs.recommended,
  importX.flatConfigs.recommended,
  importX.flatConfigs.typescript,
  
  {
    files: ["**/*.{js,mjs,cjs,ts,jsx,tsx}"],
    plugins: {
      depend,
    },
    extends: ["depend/flat/recommended"],
    languageOptions:{
      parser: tsParser,
      ecmaVersion: "latest",
      globals: globals.browser,
      parserOptions: {
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
        ecmaFeatures: {
          jsx: true
        }
      }
    },
    settings: {
      react: {
        version: "detect"
      }
    },
    rules: {
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
      "no-console": ["warn"],
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
      "react/self-closing-comp": ["error"],
      "react/jsx-wrap-multilines": ["error"],
      "react/react-in-jsx-scope": ["off"],
      "react/no-array-index-key": ["error"],
      "react/jsx-props-no-spreading": ["error"],
      "react-hooks/react-compiler": ["error"],
      "react-refresh/only-export-components": ["off"],
      "@typescript-eslint/no-unnecessary-condition": ["off"],
    }
  },

  { 
    ignores: [
      "build",
      ".react-router",
      "eslint.config.js",
      "playwright.config.ts",
      "vite.config.ts",
      "app.spec.ts",
      "vite-env.d.ts"
    ] 
  },
]);
