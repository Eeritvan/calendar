import globals from "globals";
import expoConfig from 'eslint-config-expo/flat.js';
import pluginQuery from '@tanstack/eslint-plugin-query';
import { defineConfig } from "eslint/config";
import { baseConfig } from "@repo/eslint-config/base";

export default defineConfig([
  expoConfig,
  baseConfig,
  pluginQuery.configs['flat/recommended'],
  {
    files: ['babel.config.js'],
    languageOptions: {
      globals: globals.browser,
    },
  },
  {
    ignores: [
      "dist/*",
      "android/*",
      "ios/*",
      "eslint.config.js",
      "expo-env.d.ts",
    ]
  },
]);
